package utils

import (
	"bufio"
	"io"
	"log"
	"regexp"
	"strings"
	"time"
)

type FFmpegProgress struct {
	Time  string
	Size  string
	Speed string
}

type FFmpegProgressMonitor struct {
	timeRegex        *regexp.Regexp
	speedRegex       *regexp.Regexp
	sizeRegex        *regexp.Regexp
	logInterval      time.Duration
	lastLogTime      time.Time
	progressCallback func(FFmpegProgress)
}

func NewFFmpegProgressMonitor() *FFmpegProgressMonitor {
	return &FFmpegProgressMonitor{
		timeRegex:   regexp.MustCompile(`time=(\d+:\d+:\d+\.\d+)`),
		speedRegex:  regexp.MustCompile(`speed=([\d\.]+)x`),
		sizeRegex:   regexp.MustCompile(`size=\s*([\d\.]+)kB`),
		logInterval: 2 * time.Second,
		lastLogTime: time.Now(),
	}
}

func (m *FFmpegProgressMonitor) SetLogInterval(interval time.Duration) {
	m.logInterval = interval
}

func (m *FFmpegProgressMonitor) SetProgressCallback(callback func(FFmpegProgress)) {
	m.progressCallback = callback
}

func (m *FFmpegProgressMonitor) MonitorProgress(stderr io.ReadCloser) error {
	defer stderr.Close()

	scanner := bufio.NewScanner(stderr)

	for scanner.Scan() {
		line := scanner.Text()

		if time.Since(m.lastLogTime) >= m.logInterval {
			progress := m.parseProgress(line)

			if progress.hasData() {
				if m.progressCallback != nil {
					m.progressCallback(progress)
				} else {
					m.logProgress(progress)
				}
				m.lastLogTime = time.Now()
			}
		}
	}

	return scanner.Err()
}

func (m *FFmpegProgressMonitor) parseProgress(line string) FFmpegProgress {
	progress := FFmpegProgress{}

	if timeMatch := m.timeRegex.FindStringSubmatch(line); len(timeMatch) > 1 {
		progress.Time = timeMatch[1]
	}

	if speedMatch := m.speedRegex.FindStringSubmatch(line); len(speedMatch) > 1 {
		progress.Speed = speedMatch[1] + "x"
	}

	if sizeMatch := m.sizeRegex.FindStringSubmatch(line); len(sizeMatch) > 1 {
		progress.Size = sizeMatch[1] + "kB"
	}

	return progress
}

func (p FFmpegProgress) hasData() bool {
	return p.Time != "" || p.Size != "" || p.Speed != ""
}

func (m *FFmpegProgressMonitor) logProgress(progress FFmpegProgress) {
	var progressParts []string

	if progress.Size != "" {
		progressParts = append(progressParts, "Size: "+progress.Size)
	}
	if progress.Time != "" {
		progressParts = append(progressParts, "Time: "+progress.Time)
	}
	if progress.Speed != "" {
		progressParts = append(progressParts, "Speed: "+progress.Speed)
	}

	if len(progressParts) > 0 {
		log.Println(strings.Join(progressParts, " | "))
	}
}

func (p FFmpegProgress) String() string {
	var parts []string

	if p.Size != "" {
		parts = append(parts, "Size: "+p.Size)
	}
	if p.Time != "" {
		parts = append(parts, "Time: "+p.Time)
	}
	if p.Speed != "" {
		parts = append(parts, "Speed: "+p.Speed)
	}

	return strings.Join(parts, " | ")
}

func FFmpegProgressHandler(stderr io.ReadCloser, logInterval time.Duration) error {
	monitor := NewFFmpegProgressMonitor()
	monitor.SetLogInterval(logInterval)
	return monitor.MonitorProgress(stderr)
}

func FFmpegProgressHandlerWithCallback(stderr io.ReadCloser, logInterval time.Duration, callback func(FFmpegProgress)) error {
	monitor := NewFFmpegProgressMonitor()
	monitor.SetLogInterval(logInterval)
	monitor.SetProgressCallback(callback)
	return monitor.MonitorProgress(stderr)
}
