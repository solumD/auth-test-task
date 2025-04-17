package email

import (
	"fmt"

	"github.com/solumD/auth-test-task/internal/logger"
	"github.com/solumD/auth-test-task/internal/service"
)

type emailService struct {
}

// New returns new email service object
func New() service.EmailService {
	return &emailService{}
}

// SendEmail sends an email (logs a warning message in console)
func (es *emailService) SendEmail(from string, to string, theme string, text string) {
	logger.Warn(fmt.Sprintf("Send email from %s to %s with theme %s and this text: %s", from, to, theme, text))
}
