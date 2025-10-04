package server

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/rand/v2"

	pb "github.com/Coosis/go-ecommerce/auth/internal/pb/v1/auth"
	internal "github.com/Coosis/go-ecommerce/auth/internal"
	mailmodel "github.com/Coosis/go-aws-mail-service/model"
	pgtype "github.com/jackc/pgx/v5/pgtype"
	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
	log "github.com/sirupsen/logrus"
	sqlc "github.com/Coosis/go-ecommerce/auth/sqlc"
)

// this type exists because the process of putting into valkey 
// and getting k-v pair out of valkey is decoupled and we want to retain information
type Challenge struct {
	ID      string `json:"id"`
	Digits  string `json:"digits"`
	UserID  pgtype.UUID `json:"user_id"`

	Type string `json:"type"`
	Value string `json:"value"`
}

func(c *Challenge) Marshal() (string, error) {
	data, err := json.Marshal(c)
	if err != nil {
		return "", fmt.Errorf("failed to marshal challenge: %w", err)
	}
	return base64.RawStdEncoding.EncodeToString(data), nil
}

func(c *Challenge) Unmarshal(data string) error {
	decoded, err := base64.RawStdEncoding.DecodeString(data)
	if err != nil {
		return fmt.Errorf("failed to decode challenge: %w", err)
	}
	if err := json.Unmarshal(decoded, c); err != nil {
		return fmt.Errorf("failed to unmarshal challenge: %w", err)
	}
	return nil
}

func randFourDigit() string {
	var digits = "0123456789"
	var result string
	for range 4 {
		result += string(digits[rand.IntN(len(digits))])
	}
	return result
}

func (s *Server) AnswerChallenge(
	ctx context.Context,
	req *pb.AnswerChallengeRequest,
) (*pb.Session, error) {
	ok, user_id, err := s.VerifyChallenge(ctx, req.GetChallenge(), req.GetAnswer())
	if err != nil {
		return nil, fmt.Errorf("failed to verify challenge: %v", err)
	}

	if !ok {
		return nil, fmt.Errorf("wrong code")
	}

	sessionID, err := s.CreateSession(user_id)
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %v", err)
	}
	return &pb.Session{
		SessionId: sessionID,
	}, nil
}

// 1. create a challenge id
// 2. set challenge id in valkey
// 3. send challenge id to the queue
func(s *Server) CreateChallenge(
	ctx context.Context,
	userID string,
	to string,
) (string, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return "", fmt.Errorf("creating uuid v7 failed: %v", err.Error())
	}
	challengeID := id.String()
	digits := randFourDigit()

	ch, err := s.AmqpClient.Channel()
	if err != nil {
		return "", fmt.Errorf("failed to create amqp channel: %w", err)
	}
	defer ch.Close()

	q, err := ch.QueueDeclare(
		internal.AmqpQueue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return "", fmt.Errorf("failed to declare queue: %w", err)
	}

	Job := &mailmodel.MailJob{
		To: to,
		Subject: "Your Verification Code",
		Message: fmt.Sprintf("Your verification code is: %s", digits),
	}
	data, err := Job.Marshal()
	if err != nil {
		return "", fmt.Errorf("failed to marshal job: %w", err)
	}

	err = ch.PublishWithContext(ctx,
		"",
		q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        data,
		},
	)
	if err != nil {
		return "", fmt.Errorf("failed to publish message to queue: %w", err)
	}

	var pg_user_id pgtype.UUID
	if err := pg_user_id.Scan(userID); err != nil {
		return "", fmt.Errorf("failed to set user id: %v", err)
	}

	challenge := &Challenge{
		ID:     challengeID,
		Digits: digits,
		UserID: pg_user_id,

		Type: "email", // assuming email for now, can be extended to other types
		Value: to,
	}
	marshaled, err := challenge.Marshal()
	if err != nil {
		return "", err
	}
	log.Infof("Storing challenge in valkey: %s", marshaled)
	r := s.VKclient.Do(
		ctx,
		s.VKclient.B().Set().
		Key(fmt.Sprintf("challenge:%s", challengeID)).
		Value(marshaled).
		Build(),
	)
	if r.Error() != nil {
		return "", fmt.Errorf("failed to set challenge in valkey: %w", r.Error())
	}
	return challengeID, nil
}

func(s *Server) VerifyChallenge(
	ctx context.Context,
	challengeID string,
	digits string,
) (bool, pgtype.UUID, error) {
	val := s.VKclient.Do(
		ctx,
		s.VKclient.B().Get().
		Key(fmt.Sprintf("challenge:%s", challengeID)).
		Build(),
	)
	if val.Error() != nil {
		return false, pgtype.UUID{}, fmt.Errorf("failed to get challenge: %w", val.Error())
	}
	valstr, err := val.ToString()
	log.Infof("Verifying challenge value: %s", valstr)
	if err != nil {
		return false, pgtype.UUID{}, fmt.Errorf("failed to convert challenge value to string: %w", err)
	}
	challenge := &Challenge{}
	if err := challenge.Unmarshal(valstr); err != nil {
		return false, pgtype.UUID{}, fmt.Errorf("failed to unmarshal challenge: %w", err)
	}
	if challenge.Digits != digits {
		return false, pgtype.UUID{}, fmt.Errorf("challenge digits do not match")
	}

	queries := sqlc.New(s.Pool)
	queries.VerifyContactMethod(ctx, sqlc.VerifyContactMethodParams{
		UserID: challenge.UserID,
		Type: challenge.Type,
		Value: challenge.Value,
	})

	// delete challenge from valkey
	del := s.VKclient.Do(
		ctx,
		s.VKclient.B().Del().
		Key(fmt.Sprintf("challenge:%s", challengeID)).
		Build(),
	)
	if del.Error() != nil {
		log.Errorf("failed to delete challenge: %v", del.Error())
	}
	return true, challenge.UserID, nil
}
