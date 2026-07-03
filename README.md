# dan

`dan`은 Go로 구현된 경량 웹 프레임워크입니다. 최소한의 구조로 라우팅, 그룹 라우팅, 미들웨어, JSON 응답 및 에러 처리를 제공합니다.

## 주요 기능

- `Engine` 기반 단순 HTTP 라우팅
- `RouterGroup`을 이용한 경로 접두사 및 그룹 라우팅
- 동적 라우트 지원 (`/users/:id`)
- HTTP 메서드 라우팅 지원 (`GET`, `POST`, `PUT`, `PATCH`, `DELETE`, `OPTIONS`)
- 쿼리 파라미터 및 폼 파라미터 헬퍼 제공
- 미들웨어 체인 지원
- JSON 응답 및 에러 핸들링 편의 메서드
- 기본 로깅 미들웨어 제공

## 추가 예정 기능

- [x] 동적 라우트 지원 (`/users/:id`)
- [x] 더 많은 HTTP 메서드 지원 (`PUT`, `PATCH`, `DELETE`, `OPTIONS`)
- [x] 쿼리 파라미터 및 폼 파라미터 헬퍼 추가
- [ ] Panic Recovery 미들웨어 추가

## 설치

```bash
go get github.com/hooneun/dan
```

## 빠른 시작

```go
package main

import (
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"
)

func main() {
    app := NewEngine()

    // 기본 로거 미들웨어 등록
    app.Use(Logger())

    // 단일 엔드포인트 등록
    app.GET("/health", func(c *Context) error {
        return c.JSON(http.StatusOK, map[string]string{"status": "UP"})
    })

    // 동적 라우트 등록
    app.GET("/users/:id", func(c *Context) error {
        return c.JSON(http.StatusOK, map[string]string{"id": c.Param("id")})
    })

    // 그룹 라우팅 사용
    api := app.Group("/api/v1")
    api.GET("/users", func(c *Context) error {
        return c.JSON(http.StatusOK, map[string]string{"users": "test"})
    })

    srv := &http.Server{
        Addr:         ":8080",
        Handler:      app,
        ReadTimeout:  10 * time.Second,
        WriteTimeout: 10 * time.Second,
        IdleTimeout:  120 * time.Second,
    }

    go func() {
        log.Println("[START] Framework Server is running on :8080...")
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            log.Fatalf("[FATAL] Error starting server: %v", err)
        }
    }()

    shutdown := make(chan os.Signal, 1)
    signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

    <-shutdown
    log.Println("[SHUTDOWN] Graceful shutdown...")

    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    if err := srv.Shutdown(ctx); err != nil {
        log.Fatalf("[FATAL] Shutdown failed: %v", err)
    }

    log.Println("[STOP] Server cleanly stopped")
}
```

## 주요 타입

- `Engine`: HTTP 요청을 받아 처리하는 메인 엔진
- `RouterGroup`: 경로 접두사와 미들웨어를 묶어 관리하는 그룹
- `Context`: HTTP 요청과 응답을 간편하게 처리하는 컨텍스트 객체
- `HandlerFunc`: `*Context`를 받아 에러를 반환하는 핸들러 타입
- `MiddlewareFunc`: 다른 핸들러를 감싸는 미들웨어 타입

## 미들웨어

`dan`은 미들웨어 체인을 지원합니다. 기본적으로 `Logger()`가 등록되며, `Use`를 통해 추가 미들웨어를 적용할 수 있습니다.

```go
app.Use(Logger())
```

미들웨어는 `HandlerFunc`을 받아 `HandlerFunc`을 반환하는 함수 형태입니다.

## Context 헬퍼

`Context`는 다음 메서드를 제공합니다.

- `JSON(statusCode int, data any) error`
- `Error(statusCode int, message string) error`
- `BindJSON(v any) error`
- `Param(key string) string`
- `Query(key string) string`
- `DefaultQuery(key, defaultValue string) string`
- `Form(key string) string`
- `DefaultForm(key, defaultValue string) string`

예시:

```go
return c.JSON(http.StatusOK, map[string]string{"message": "hello"})
```

동적 라우트 파라미터는 `Param`으로 가져올 수 있습니다.

```go
app.GET("/users/:id", func(c *Context) error {
    id := c.Param("id")
    return c.JSON(http.StatusOK, map[string]string{"id": id})
})
```

쿼리 파라미터는 `Query` 또는 `DefaultQuery`로 가져올 수 있습니다.

```go
app.GET("/search", func(c *Context) error {
    q := c.Query("q")
    page := c.DefaultQuery("page", "1")
    return c.JSON(http.StatusOK, map[string]string{"q": q, "page": page})
})
```

폼 파라미터는 `Form` 또는 `DefaultForm`으로 가져올 수 있습니다.

```go
app.POST("/profile", func(c *Context) error {
    name := c.Form("name")
    role := c.DefaultForm("role", "guest")
    return c.JSON(http.StatusOK, map[string]string{"name": name, "role": role})
})
```

## HTTP 메서드

다음 HTTP 메서드 라우팅을 지원합니다.

- `GET(path string, h HandlerFunc)`
- `POST(path string, h HandlerFunc)`
- `PUT(path string, h HandlerFunc)`
- `PATCH(path string, h HandlerFunc)`
- `DELETE(path string, h HandlerFunc)`
- `OPTIONS(path string, h HandlerFunc)`

예시:

```go
app.PUT("/users/:id", func(c *Context) error {
    return c.JSON(http.StatusOK, map[string]string{"id": c.Param("id")})
})
```

## 예제

이 저장소의 `example` 파일에 기본 사용 예제가 포함되어 있습니다.

## 실행 방법

1. 프로젝트 루트로 이동합니다.
2. `go run example` 또는 `go run .` 명령으로 서버를 실행합니다.
3. 브라우저 또는 HTTP 클라이언트로 `http://localhost:8080/health` 등에 접근합니다.

## 패키지 임포트

```go
import "github.com/hooneun/dan"
```

## 기여

1. 저장소를 포크합니다.
2. 새로운 브랜치를 만듭니다.
3. 수정사항을 커밋합니다.
4. 풀 리퀘스트(PR)를 생성합니다.

기능 추가, 버그 수정, 문서 개선 모두 환영합니다.

## 라이선스

별도로 명시된 라이선스가 없으므로 자유롭게 사용하실 수 있습니다.
