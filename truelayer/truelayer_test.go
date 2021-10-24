package truelayer_test

import (
	"net/url"
	"testing"

	"github.com/ImTomEddy/truelayer-go/truelayer"
	"github.com/ImTomEddy/truelayer-go/truelayer/providers"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	tl := truelayer.New("abc", "abc", false)
	require.NotNil(t, tl)
}

func TestGetAuthenticationLink(t *testing.T) {
	t.Run("sandbox link", func(t *testing.T) {
		expected := "https://auth.truelayer-sandbox.com/?client_id=abc&providers=uk-ob-monzo&redirect_uri=http%3A%2F%2Flocalhost&response_type=code&scope=balance"
		tl := truelayer.New("abc", "abc", true)
		require.NotNil(t, tl)

		uri, err := url.Parse("http://localhost")
		require.NoError(t, err)

		link, err := tl.GetAuthenticationLink([]string{providers.UKMonzo}, []string{truelayer.PermissionBalance}, uri, false)
		require.NoError(t, err)
		require.NotEmpty(t, link)
		require.Equal(t, expected, link)
	})

	t.Run("live link", func(t *testing.T) {
		expected := "https://auth.truelayer.com/?client_id=abc&providers=uk-ob-monzo&redirect_uri=http%3A%2F%2Flocalhost&response_type=code&scope=balance"
		tl := truelayer.New("abc", "abc", false)
		require.NotNil(t, tl)

		uri, err := url.Parse("http://localhost")
		require.NoError(t, err)

		link, err := tl.GetAuthenticationLink([]string{providers.UKMonzo}, []string{truelayer.PermissionBalance}, uri, false)
		require.NoError(t, err)
		require.NotEmpty(t, link)
		require.Equal(t, expected, link)
	})

	t.Run("live link post code", func(t *testing.T) {
		expected := "https://auth.truelayer.com/?client_id=abc&providers=uk-ob-monzo+uk-ob-mbna&redirect_uri=http%3A%2F%2Flocalhost&response_mode=form_post&response_type=code&scope=balance"
		tl := truelayer.New("abc", "abc", false)
		require.NotNil(t, tl)

		uri, err := url.Parse("http://localhost")
		require.NoError(t, err)

		link, err := tl.GetAuthenticationLink([]string{providers.UKMonzo, providers.UKMBNA}, []string{truelayer.PermissionBalance}, uri, true)
		require.NoError(t, err)
		require.NotEmpty(t, link)
		require.Equal(t, expected, link)
	})
}
