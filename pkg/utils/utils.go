/**
 * @Author:      leafney
 * @GitHub:      https://github.com/leafney
 * @Project:     music-grabber
 * @Date:        2024-01-13 15:12
 * @Description:
 */

package utils

import "fmt"

func FormatFileSize(body []byte) string {
	sizeInMB := float64(len(body)) / (1024 * 1024)
	if sizeInMB > 1024 {
		sizeInGB := sizeInMB / 1024
		remainMB := sizeInMB - (sizeInGB * 1024)
		return fmt.Sprintf("%.2f GB %.2f MB", sizeInGB, remainMB)
	}
	return fmt.Sprintf("%.2f MB", sizeInMB)
}
