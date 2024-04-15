package telegram

import (
	"github.com/omegaatt36/instagramrobot/domain"
	"github.com/omegaatt36/instagramrobot/logging"
	"github.com/pkg/errors"
	"gopkg.in/telebot.v3"
)

type MediaSender struct {
	bot *telebot.Bot
	msg *telebot.Message
}

// NewMediaSenderImpl creates a new MediaSenderImpl instance
func NewMediaSenderImpl(bot *telebot.Bot, msg *telebot.Message) domain.MediaSender {
	return &MediaSender{
		bot: bot,
		msg: msg,
	}
}

const (
	caption = "🤖 Downloaded via @Devigdlbot *Powered by* - @DEVBOTS2"
)

// Send will start to process Media and eventually send it to the Telegram chat
func (m *MediaSender) Send(media *domain.Media) error {
	logging.Infof("chatID(%d) shortcode(%s)", m.msg.Sender.ID, media.Shortcode)

	// Check if media has no child item
	if len(media.Items) == 0 {
		return m.sendSingleMedia(media)
	}

	return m.sendNestedMedia(media)
}

func (m *MediaSender) sendSingleMedia(media *domain.Media) error {
	if media.IsVideo {
		if _, err := m.bot.Send(m.msg.Chat, &telebot.Video{
			File:    telebot.FromURL(media.Url),
			Caption: caption,
		}); err != nil {
			return errors.Wrap(err, "couldn't send the single video")
		}

		logging.Debugf("Sent single video with short code [%v]", media.Shortcode)
	} else {
		if _, err := m.bot.Send(m.msg.Chat, &telebot.Photo{
			File:    telebot.FromURL(media.Url),
			Caption: caption,
		}); err != nil {
			return errors.Wrap(err, "couldn't send the single photo")
		}

		logging.Debugf("Sent single photo with short code [%v]", media.Shortcode)
	}

	return m.SendCaption(media)
}

func (m *MediaSender) sendNestedMedia(media *domain.Media) error {
	_, err := m.bot.SendAlbum(m.msg.Chat, m.generateAlbumFromMedia(media))
	if err != nil {
		return errors.Wrap(err, "couldn't send the nested media")
	}
	return m.SendCaption(media)
}

func (m *MediaSender) generateAlbumFromMedia(media *domain.Media) telebot.Album {
	var album telebot.Album

	for _, media := range media.Items {
		if media.IsVideo {
			album = append(album, &telebot.Video{
				File: telebot.FromURL(media.Url),
			})
		} else {
			album = append(album, &telebot.Photo{
				File: telebot.FromURL(media.Url),
			})
		}
	}

	return album
}

// SendCaption will send the caption to the chat.
func (m *MediaSender) SendCaption(media *domain.Media) error {
	// If caption is empty, ignore sending it
	if media.Caption == "" {
		return nil
	}
	// TODO: chunk caption if the length is above the Telegram limit
	_, err := m.bot.Reply(m.msg, media.Caption)
	return err
}
