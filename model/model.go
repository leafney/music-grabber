/**
 * @Author:      leafney
 * @GitHub:      https://github.com/leafney
 * @Project:     music-grabber
 * @Date:        2024-01-13 11:37
 * @Description:
 */

package model

import (
	"context"
	"github.com/chromedp/cdproto/target"
)

type TargetId struct {
	Id      target.ID
	Type    int
	BaseCtx context.Context
}

type UrlLink struct {
	Url string
	Ctx context.Context
}
