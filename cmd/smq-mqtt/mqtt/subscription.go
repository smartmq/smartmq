package mqtt

import (
	"log"
	"regexp"
	"strings"
)

type Subscription struct {
	TopicFilter string
	Qos         byte
	Regexp      *regexp.Regexp
	cache       map[string]bool
}

func NewSubscription(topic string, qos byte) *Subscription {
	re, err := toRegexPattern(topic)
	if err != nil {
		log.Fatal(err)
	}
	s := &Subscription{
		TopicFilter: topic,
		Qos:         qos,
		Regexp:      re,
		cache:       make(map[string]bool),
	}
	return s
}
func toRegexPattern(subscribedTopic string) (*regexp.Regexp, error) {
	var regexPattern string
	regexPattern = subscribedTopic
	regexPattern = strings.Replace(regexPattern, "#", ".*", -1)
	regexPattern = strings.Replace(regexPattern, "+", "[^/]*", -1)
	pattern, err := regexp.Compile("^" + regexPattern + "$")
	return pattern, err
}

func (s *Subscription) IsSubscribed(publishingTopic string) bool {
	// first get from cache ...
	var ret, present bool
	ret, present = s.cache[publishingTopic]
	if present {
		return ret
	} else {
		if strings.Compare(s.TopicFilter, publishingTopic) == 0 {
			ret = true
		} else {
			topicMatches := s.Regexp.MatchString(publishingTopic)
			ret = topicMatches
		}

		// put result of calculation in cache ...
		s.cache[publishingTopic] = ret

		return ret
	}
}
