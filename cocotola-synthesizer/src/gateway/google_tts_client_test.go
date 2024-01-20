//go:build large

package gateway_test

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/gopxl/beep/mp3"

	"github.com/kujilabo/cocotola-1.21/cocotola-synthesizer/src/gateway"
	libdomain "github.com/kujilabo/cocotola-1.21/lib/domain"
	"github.com/stretchr/testify/assert"
)

func Test_googleTTSClient_Synthesize(t *testing.T) {
	t.Skip()
	ctx := context.Background()
	t.Parallel()
	httpClient := http.Client{}
	apiKey := os.Getenv("GOOGLE_TEXT_TO_SPEECH_API_KEY")
	c := gateway.NewGoogleTTSClient(&httpClient, apiKey)
	audioContent, err := c.Synthesize(ctx, libdomain.Lang5JAJP, "FEMALE", "こんにちは")
	assert.Nil(t, err)
	fmt.Println(audioContent)
	log.Fatal(audioContent)
}

func Test_Audio(t *testing.T) {
	audioContent := "//NExAAAAANIAAAAAC4LB/i3hqxDyFnWo1ArGR5T3fq9DEMVisVjx5TKAAAAAAhNwAQnXc4HFuf/9d3c9AMH6P4nD+p0o7///gh/P+UdiAHwfeIATB8H1Ve9FuXxCtWK//NExFMSkQH0AU8YANVb1dSCCmONvrhwfZa1RAYG9QBj6lszLZQB0cL2gZMGBqQKKnUy3STg3MEFxAMfQ7N+k3UXS8OAWAi5/r/7JGJPjJkwgj/3ofmRUIgZzAojjFX///NExFsfkyo4AZigAP/6HIiLLDC45hfFjIYVCqbE4VP+r/b/+mpR03K7LNzI0PH0FYEB9uCxVRh7nWsK1vZwdLHFTxPnxQBYkVcyMhYiKCUA48dwj0AqEEScIcbDLDLl//NExC8b4iqYAY+AALTSNCeGWMzheSraXzdNF3LpiYqdf9anXbUpXVqupk90EHr65kXf61OXEpoaJzN2TPFkmXKriF52ln0/WuLWnP6V/XbsZC0IrHZbkVAC5GJSOM3y//NExBIWQgasAdiAAKnFvw/IsY5BNatTzZCCfxzTIvCugakWNlOUyTRbLiXUh6Ddafpv1G/mDdIcxF8mDR9SDbmBp1l92039X0H1m6z05////owY9gCzTSoiDAe5SbDU//NExAwUYZ7AAGyKlAEZlHMhBYdi3REuBs+MmgiZiIF9B05eNGzpun1+3/p9+ioQ9QwPgceMFxdCTDUECI4IcPkXMEaHHmi4RNAALBgx////+MX1hJQc6WkAZimaLwfY//NExA0VAarEAG4KlEIpdXyxpWQ1qvabqApSx/KKs4T9yu33HGzbz//ybsn///ZnZaOUehAGA4EEAQcOEBMgoUeJhBg0GF3UCHi5ZoUMIP/////o9RWEQWPYYwAtApB7//NExAwVQarEAGvKlHVeA1QqCwV3BL6WhluErAK8RYXCjbHUlnN5E+a4xTUzsd2X///7uVDkKNEQ4DiYkEBhRwuIGHBwzBogIHninj0NCRYqn////+n6xYhYaCzAScLh//NExAoSWa7IAGvKlK3WWpIU5NvVoSdtrekUJMLURCg1nWXs7zVPyuT2f///3ctyR4sPAodFw4RmZxiMUcgiJEAjEs+08AGf////tSr6zQLMlybu5RDAppX6IJIMn/5N//NExBMRKcbAAILElEdK/vOjcHkhH/bUQxIE9VSlJ9Wf6/2f5UzUDB1OyGcqWmdWVDKUKQWc+Wd////98KzF///XYCDblqyp9MtUzwjQ+u1JnhPjJfMieGtektcpPqWO//NExCESGxqgAMtKvdEyKZr1seQ+pvrb9JvlJ9Xf9Po6/KKzf/R//////5v7//7+dACDgCrv7wzlCS5rMYkjaW1uxXoHrC4Kl9nZjETi+omgsCGr0Spvub/UfPfUs0/Q//NExCsSOyKUANNKvP1/lK/Xf+//8y//2//////UpRf//RvKYBzyh0Oiz7WOW7klUWMcrC10pXTIrPXQCEV7M01O8VHqXArGCfqIq31n/0W/T/Uf/R/U3/+39f7//7////NExDURYxqMANwKvf//0+iCLkrz/26UFAaRxZqV7+GdiXtMDSQtsRsW1N63DQKdz/1pLHwUW2mIhwhE/UOJvrV+cNP0X/T//G/0/v+hP7/r//////k9DDD/p/V88weA//NExEIROxKMANNKuXMKKL+t6u0CbopZpG34xFb96Vg4ghXFLWxOCRm/UYjFAcBp8xIH62/Mf7/2/Rvqaf+v//b//T//////UxrUHix7/f+Tsu4WEIj3H39ApgQ+CMZD//NExFAR+xaMANtOuckdZwMXDszODxG9hH/m//9/1X9Ru3xQDghXzFfyoOy7KZccEhvoy/U/8cM//////qe//9DDG6j557zEKExyf0DMID4MVi1lVV2GeDTXnFVBiLPz//NExFsRaxqoAJnOvRP//6/55/yr8wGQe1OiJ+CKQTBTMpAkWZ7qWF0I5+7LHwQMoKqmA4AsX16Zs+umSBm37frNvrTb61t+dS+pB/aTkemyBIl72GDVf2Ng+gPofZZS//NExGgXmwagAKNauVqKYYHJHvUG8t9R8//+/5j/MJCX1Ix83jAFgALMQ3YdTcwwQAFElqxgqDUBom3MH4yFVnx+TljvZiN/0/o+Y1Df7fqc/MdP7kptCQgMFcWjlnnu//NExFwXox6sAJKUvGD41X6iyGzgZ5B9jppTUH/IEnyBAE3WhW//3//VDfdVqoosgPiioNQXHB8LRTkFB+CkcQtih4NgbC0ykoLGjrhRYVRqv6hm/jXr//+4hjVXuP////NExFAXkuKsAJiQua/fn4qfvquu1GSyppss1f/4QywACBKqiR2kWIbyFB9nQSOxrgQBJRRFvcj30ZP09t/bVaPhRJSzASHn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"

	audioContentBytes, err := base64.StdEncoding.DecodeString(audioContent)
	assert.NoError(t, err)

	reader := bytes.NewReader(audioContentBytes)
	readCloser := io.NopCloser(reader)
	// r := io.NopCloser(buffer.New()(bytes))

	streamer, format, err := mp3.Decode(readCloser)
	if err != nil {
		log.Fatal(err)
	}
	defer streamer.Close()
	fmt.Println(format.SampleRate)

}
