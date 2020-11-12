package fake

import (
	"testing"

	"github.com/xh3b4sd/logger"
)

func Test_Fake_Interface(t *testing.T) {
	var _ logger.Interface = New()
}
