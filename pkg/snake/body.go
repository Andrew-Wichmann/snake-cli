package snake

import (
	"fmt"
	"log"
	"math"
)

type body struct {
	max_length float64
	segments   []segment
	growRate   float64
}

func newBody() body {
	head := point{50, 50}
	tail := point{50, 50}
	max_length := 100.0
	segments := []segment{{&tail, &head, UP}}
	return body{segments: segments, growRate: 1.0, max_length: max_length}
}

func directionChangeLegal(current, requested Direction) bool {
	if current == requested {
		return false
	}
	if current == UP || current == DOWN {
		return requested == RIGHT || requested == LEFT
	} else if current == RIGHT || current == LEFT {
		return requested == UP || requested == DOWN
	} else {
		panic(fmt.Sprintf("Invalid direction change current=%d requested=%d", current, requested))
	}
}

func (b *body) changeDirection(direction Direction) {
	head := b.segments[len(b.segments)-1]
	if directionChangeLegal(head.orientation, direction) {
		newHead := point{head.end.x, head.end.y}
		segments := append(b.segments, segment{head.end, &newHead, direction})
		b.segments = segments
	}
}

func (b *body) grow() {
	growingSegment := b.segments[len(b.segments)-1]
	if growingSegment.orientation == UP {
		growingSegment.end.y -= b.growRate
	}
	if growingSegment.orientation == DOWN {
		growingSegment.end.y += b.growRate
	}
	if growingSegment.orientation == RIGHT {
		growingSegment.end.x += b.growRate
	}
	if growingSegment.orientation == LEFT {
		growingSegment.end.x -= b.growRate
	}
	b.segments[len(b.segments)-1] = growingSegment
	l := b.length()
	if l > b.max_length {
		b.shorten(l - b.max_length)
	}
}

func (b *body) length() float64 {
	var l float64
	for _, segment := range b.segments {
		l += segment.length()
	}
	return l
}

func (b *body) shorten(amount float64) {
	for i := 0; len(b.segments) > 0; {
		amount := b.segments[i].shorten(amount)
		if amount <= 0 {
			return
		}
		b.segments = b.segments[1:]
	}
}

type segment struct {
	start       *point
	end         *point
	orientation Direction
}

func (s segment) length() float64 {
	a2 := math.Pow(math.Abs(s.start.x-s.end.x), 2)
	b2 := math.Pow(math.Abs(s.start.y-s.end.y), 2)
	c := math.Sqrt(a2 + b2)
	log.Printf("a2=%f b2=%f c=%f\n", a2, b2, c)
	return c
}

func (s *segment) shorten(amount float64) float64 {
	l := s.length()
	log.Printf("%f\n", l)
	if s.orientation == UP {
		s.start.y -= amount
	}
	if s.orientation == RIGHT {
		s.start.x += amount
	}
	if s.orientation == DOWN {
		s.start.y += amount
	}
	if s.orientation == LEFT {
		s.start.x -= amount
	}
	return amount - l
}

type point struct {
	x, y float64
}
