package core

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.dagger.io/dagger/core"
	"go.dagger.io/dagger/internal/testutil"
)

func TestSecretEnvFromFile(t *testing.T) {
	t.Parallel()

	secretID := secretID(t, "some-content")

	var envRes struct {
		Container struct {
			From struct {
				WithSecretVariable struct {
					Exec struct {
						Stdout struct{ Contents string }
					}
				}
			}
		}
	}

	err := testutil.Query(
		`query Test($secret: SecretID!) {
			container {
				from(address: "alpine:3.16.2") {
					withSecretVariable(name: "SECRET", secret: $secret) {
						exec(args: ["env"]) {
							stdout { contents }
						}
					}
				}
			}
		}`, &envRes, &testutil.QueryOptions{Variables: map[string]any{
			"secret": secretID,
		}})
	require.NoError(t, err)
	require.Contains(t, envRes.Container.From.WithSecretVariable.Exec.Stdout.Contents, "SECRET=some-content\n")
}

func TestSecretMountFromFile(t *testing.T) {
	t.Parallel()

	secretID := secretID(t, "some-content")

	var envRes struct {
		Container struct {
			From struct {
				WithMountedSecret struct {
					Exec struct {
						Stdout struct{ Contents string }
					}
				}
			}
		}
	}

	err := testutil.Query(
		`query Test($secret: SecretID!) {
			container {
				from(address: "alpine:3.16.2") {
					withMountedSecret(path: "/sekret", source: $secret) {
						exec(args: ["cat", "/sekret"]) {
							stdout { contents }
						}
					}
				}
			}
		}`, &envRes, &testutil.QueryOptions{Variables: map[string]any{
			"secret": secretID,
		}})
	require.NoError(t, err)
	require.Contains(t, envRes.Container.From.WithMountedSecret.Exec.Stdout.Contents, "some-content")
}

func secretID(t *testing.T, content string) core.SecretID {
	var secretRes struct {
		Directory struct {
			WithNewFile struct {
				File struct {
					Secret struct {
						ID core.SecretID
					}
				}
			}
		}
	}

	err := testutil.Query(
		`query Test($content: String!) {
			directory {
				withNewFile(path: "some-file", contents: $content) {
					file(path: "some-file") {
						secret {
							id
						}
					}
				}
			}
		}`, &secretRes, &testutil.QueryOptions{Variables: map[string]any{
			"content": content,
		}})
	require.NoError(t, err)

	secretID := secretRes.Directory.WithNewFile.File.Secret.ID
	require.NotEmpty(t, secretID)

	return secretID
}