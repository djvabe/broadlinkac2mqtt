package webClient

import (
	"context"
	"log/slog"
	"net"
	"strconv"
	"time"

	"github.com/djvabe/broadlinkac2mqtt/app"
	"github.com/djvabe/broadlinkac2mqtt/app/webClient/models"
)

type webClient struct {
	logger *slog.Logger
}

func NewWebClient(logger *slog.Logger) app.WebClient {
	return &webClient{
		logger: logger,
	}
}

func (w *webClient) SendCommand(ctx context.Context, input *models.SendCommandInput) (*models.SendCommandReturn, error) {
	var (
		err      error
		conn     net.Conn
		response []byte
	)

	for i := 0; i < 3; i++ {
		response, err = w.sendCommand(ctx, input)
		if err == nil {
			return &models.SendCommandReturn{Payload: response}, nil
		}
		w.logger.WarnContext(ctx, "Failed to send command, retrying...", slog.Int("attempt", i+1), slog.Any("err", err))
		time.Sleep(time.Second * 1)
	}

	return nil, err
}

func (w *webClient) sendCommand(ctx context.Context, input *models.SendCommandInput) ([]byte, error) {
	conn, err := net.DialTimeout("udp", input.Ip+":"+strconv.Itoa(int(input.Port)), time.Second*5)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	err = conn.SetDeadline(time.Now().Add(time.Second * 5))
	if err != nil {
		return nil, err
	}

	_, err = conn.Write(input.Payload)
	if err != nil {
		return nil, err
	}

	response := make([]byte, 1024)
	n, err := conn.Read(response)
	if err != nil {
		return nil, err
	}

	return response[:n], nil
}
