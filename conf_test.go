package structconf_test

import (
	"context"
	"flag"
	"testing"

	"github.com/kevinfalting/structconf"
	"github.com/kevinfalting/structconf/stronf"
)

func TestParse(t *testing.T) {
	t.Run("mw & handlers", func(t *testing.T) {
		type thing struct {
			Int int
		}

		t.Run("no middleware no handlers", func(t *testing.T) {
			a := thing{Int: 22}

			conf := structconf.Conf[thing]{}
			if err := conf.Parse(context.Background(), &a); err != nil {
				t.Error(err)
			}

			if a.Int != 22 {
				t.Errorf("expected 22, got %d", a.Int)
			}
		})

		t.Run("one handler, no middleware", func(t *testing.T) {
			a := thing{}

			conf := structconf.Conf[thing]{
				Handlers: []stronf.Handler{
					stronf.HandlerFunc(
						func(ctx context.Context, f stronf.Field, _ any) (any, error) {
							return int(5), nil
						},
					),
				},
			}

			if err := conf.Parse(context.Background(), &a); err != nil {
				t.Error(err)
			}

			if a.Int != 5 {
				t.Errorf("expected 5, got %d", a.Int)
			}
		})
		t.Run("one middleware, no handler", func(t *testing.T) {
			a := thing{}

			conf := structconf.Conf[thing]{
				Middleware: []stronf.Middleware{
					func(h stronf.Handler) stronf.Handler {
						return stronf.HandlerFunc(
							func(ctx context.Context, f stronf.Field, iv any) (any, error) {
								return h.Handle(ctx, f, iv)
							},
						)
					},
				},
			}

			if err := conf.Parse(context.Background(), &a); err != nil {
				t.Error(err)
			}
		})
	})

	t.Run("default Conf", func(t *testing.T) {
		ctx := context.Background()

		t.Run("handler overwrites val in struct", func(t *testing.T) {
			type A struct {
				Int int `conf:"env:NAME"`
			}

			conf, err := structconf.New[A]()
			if err != nil {
				t.Fatalf("failed to get a new conf: %v", err)
			}

			t.Setenv("NAME", "22")

			a := A{
				Int: 88,
			}

			if err := conf.Parse(ctx, &a); err != nil {
				t.Errorf("expected no error, got %v", err)
			}

			if a.Int != 22 {
				t.Errorf("expected 22, got %d", a.Int)
			}
		})

		t.Run("handler returns nil, does not overwrite val in struct", func(t *testing.T) {
			type A struct {
				Int int `conf:"env:NAME"`
			}

			conf, err := structconf.New[A]()
			if err != nil {
				t.Fatalf("failed to get a new conf: %v", err)
			}

			a := A{
				Int: 88,
			}

			if err := conf.Parse(ctx, &a); err != nil {
				t.Errorf("expected no error, got %v", err)
			}

			if a.Int != 88 {
				t.Errorf("expected 88, got %d", a.Int)
			}
		})

		t.Run("required", func(t *testing.T) {
			type A struct {
				Int int `conf:"env:NAME,required"`
			}

			conf, err := structconf.New[A]()
			if err != nil {
				t.Fatalf("failed to get a new conf: %v", err)
			}

			t.Run("set env val", func(t *testing.T) {
				t.Setenv("NAME", "5")

				var a A

				err := conf.Parse(ctx, &a)
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}

				if a.Int != 5 {
					t.Errorf("expected 5, got %d", a.Int)
				}
			})

			t.Run("no env val set but struct has a val set", func(t *testing.T) {
				a := A{
					Int: 8,
				}

				err := conf.Parse(ctx, &a)
				if err != nil {
					t.Errorf("expected no error, got %v", err)
				}

				if a.Int != 8 {
					t.Errorf("expected 8, got %d", a.Int)
				}
			})
		})

		t.Run("env & flag using the default tag", func(t *testing.T) {
			type A struct {
				Int int `conf:"env:INT,flag:int,default:55"`
			}

			t.Run("environment variable set", func(t *testing.T) {
				var a A

				t.Setenv("INT", "22")

				fset := flag.NewFlagSet(t.Name(), flag.ContinueOnError)
				conf, err := structconf.New[A](structconf.WithFlagSet(fset))
				if err != nil {
					t.Fatalf("failed to get a new conf: %v", err)
				}

				if err := fset.Parse(nil); err != nil {
					t.Fatalf("failed to Parse flagset: %v", err)
				}

				if err := conf.Parse(context.Background(), &a); err != nil {
					t.Errorf("expected no error, got %v", err)
				}

				if a.Int != 22 {
					t.Errorf("expected 22, got %d", a.Int)
				}
			})
			t.Run("flag set", func(t *testing.T) {
				var a A

				fset := flag.NewFlagSet(t.Name(), flag.ContinueOnError)
				conf, err := structconf.New[A](structconf.WithFlagSet(fset))
				if err != nil {
					t.Fatalf("failed to get a new conf: %v", err)
				}

				if err := fset.Parse([]string{"-int", "22"}); err != nil {
					t.Fatalf("failed to Parse flagset: %v", err)
				}

				if err := conf.Parse(context.Background(), &a); err != nil {
					t.Errorf("expected no error, got %v", err)
				}

				if a.Int != 22 {
					t.Errorf("expected 22, got %d", a.Int)
				}
			})
			t.Run("flag & env set, expect flag to overwrite", func(t *testing.T) {
				var a A

				t.Setenv("INT", "456")

				fset := flag.NewFlagSet(t.Name(), flag.ContinueOnError)
				conf, err := structconf.New[A](structconf.WithFlagSet(fset))
				if err != nil {
					t.Fatalf("failed to get a new conf: %v", err)
				}

				if err := fset.Parse([]string{"-int", "22"}); err != nil {
					t.Fatalf("failed to Parse flagset: %v", err)
				}

				if err := conf.Parse(context.Background(), &a); err != nil {
					t.Errorf("expected no error, got %v", err)
				}

				if a.Int != 22 {
					t.Errorf("expected 22, got %d", a.Int)
				}
			})
			t.Run("no flag or env set, expect default to be set", func(t *testing.T) {
				var a A

				fset := flag.NewFlagSet(t.Name(), flag.ContinueOnError)
				conf, err := structconf.New[A](structconf.WithFlagSet(fset))
				if err != nil {
					t.Fatalf("failed to get a new conf: %v", err)
				}

				if err := fset.Parse(nil); err != nil {
					t.Fatalf("failed to Parse flagset: %v", err)
				}

				if err := conf.Parse(context.Background(), &a); err != nil {
					t.Errorf("expected no error, got %v", err)
				}

				if a.Int != 55 {
					t.Errorf("expected 55, got %d", a.Int)
				}
			})
		})
	})
}
