/*
 * Copyright (c) 2021 Accurics - All Rights Reserved
 * Unauthorized copying of this file, via any medium is strictly prohibited
 * Proprietary and confidential
 */

package service_test

import (
	"context"
	"testing"

	"github.com/girishg4t/app_invite_service/pkg/middleware"
	"github.com/girishg4t/app_invite_service/pkg/model"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func StubToken(token model.Token) context.Context {
	return context.WithValue(context.Background(), middleware.AuthCtxKey, token)
}

func TestService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Services Suite")
}
