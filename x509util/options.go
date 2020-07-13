package x509util

import (
	"bytes"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/pkg/errors"
	"github.com/smallstep/cli/config"
)

// Options are the options that can be passed to NewCertificate.
type Options struct {
	CertBuffer *bytes.Buffer
}

func (o *Options) apply(cr *x509.CertificateRequest, opts []Option) (*Options, error) {
	for _, fn := range opts {
		if err := fn(cr, o); err != nil {
			return o, err
		}
	}
	return o, nil
}

// Option is the type used as a variadic argument in NewCertificate.
type Option func(cr *x509.CertificateRequest, o *Options) error

// WithTemplate is an options that executes the given template text with the
// given data.
func WithTemplate(text string, data TemplateData) Option {
	return func(cr *x509.CertificateRequest, o *Options) error {
		tmpl, err := template.New("template").Funcs(sprig.TxtFuncMap()).Parse(text)
		if err != nil {
			return errors.Wrapf(err, "error parsing template")
		}

		buf := new(bytes.Buffer)
		data.SetCertificateRequest(cr)
		if err := tmpl.Execute(buf, data); err != nil {
			return errors.Wrapf(err, "error executing template")
		}
		fmt.Println(buf.String())
		o.CertBuffer = buf
		return nil
	}
}

// WithTemplateFile is an options that reads the template file and executes it
// with the given data.
func WithTemplateFile(path string, data TemplateData) Option {
	return func(cr *x509.CertificateRequest, o *Options) error {
		filename := config.StepAbs(path)
		b, err := ioutil.ReadFile(filename)
		if err != nil {
			return errors.Wrapf(err, "error reading %s", path)
		}

		tmpl, err := template.New(path).Funcs(sprig.TxtFuncMap()).Parse(string(b))
		if err != nil {
			return errors.Wrapf(err, "error parsing %s", path)
		}

		buf := new(bytes.Buffer)
		data.SetCertificateRequest(cr)
		if err := tmpl.Execute(buf, data); err != nil {
			return errors.Wrapf(err, "error executing %s", path)
		}
		o.CertBuffer = buf
		return nil
	}
}
