package tool

import (
	"fmt"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
)

type Bundle struct {
	bundle *i18n.Bundle
}

func NewBundle(lang language.Tag) *Bundle {
	b := Bundle{bundle: i18n.NewBundle(lang)}
	b.bundle.RegisterUnmarshalFunc("yaml", yaml.Unmarshal)
	return &b
}

func (b Bundle) GetMsgByCode(lang string, code int) string {
	localize := i18n.NewLocalizer(b.bundle, lang)
	msg, err := localize.Localize(&i18n.LocalizeConfig{
		MessageID: fmt.Sprintf("%d", code),
	})
	if nil != err {
		Logger.Error(err.Error())
	}
	return msg
}

func (b Bundle) GetMsgById(lang string, id string) string {
	localize := i18n.NewLocalizer(b.bundle, lang)
	msg, err := localize.Localize(&i18n.LocalizeConfig{
		MessageID: id,
	})
	if nil != err {
		Logger.Error(err.Error())
	}
	return msg
}

func (b Bundle) GetByConfig(lang string, config *i18n.LocalizeConfig) string {
	localize := i18n.NewLocalizer(b.bundle, lang)
	msg, err := localize.Localize(config)
	if nil != err {
		Logger.Error(err.Error())
	}
	return msg
}

func (b Bundle) MustLoadMessageFile(path string) {
	b.bundle.MustLoadMessageFile(path)
}
