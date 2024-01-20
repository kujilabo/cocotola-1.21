package handler

import (
	"github.com/gin-gonic/gin"
)

type SynthesizeParameter struct {
	Lang5 string `json:"lang2" binding:"required,len=5"`
	Text  string `json:"text" binding:"required"`
}

type AudioResponse struct {
	ID      int    `json:"id"`
	Lang5   string `json:"lang5"`
	Text    string `json:"text"`
	Content string `json:"content"`
}

type SynthesizerInterface interface {
}

type SynthesizerHandler struct {
	syntheziserUsecase SynthesizerInterface
}

func NewSynthesizerHandler(syntheziserUsecase SynthesizerInterface) *SynthesizerHandler {
	return &SynthesizerHandler{
		syntheziserUsecase: syntheziserUsecase,
	}
}

func (h *SynthesizerHandler) Synthesize(c *gin.Context) {
}

func (h *SynthesizerHandler) FindAudioByID(c *gin.Context) {

}

// func (h *SynthesizerHandler) errorHandle(ctx context.Context, logger *slog.Logger, c *gin.Context, err error) bool {
// 	// if errors.Is(err, service.ErrAudioNotFound) {
// 	// 	logger.Warnf("PrivateSynthesizerHandler err: %+v", err)
// 	// 	c.JSON(http.StatusNotFound, gin.H{"message": "Audio not found"})
// 	// 	return true
// 	// }
// 	logger.ErrorContext(ctx, fmt.Sprintf("SynthesizerHandler. error: %+v", err))
// 	return false
// }

func NewInitSynthesizerRouterFunc(syntheziserUsecase SynthesizerInterface) InitRouterGroupFunc {
	return func(parentRouterGroup *gin.RouterGroup, middleware ...gin.HandlerFunc) error {
		workbook := parentRouterGroup.Group("synthesize")
		SynthesizerHandler := NewSynthesizerHandler(syntheziserUsecase)
		for _, m := range middleware {
			workbook.Use(m)
		}
		workbook.POST("synthesize", SynthesizerHandler.Synthesize)
		workbook.GET("audio/:audioID", SynthesizerHandler.FindAudioByID)
		return nil
	}
}
