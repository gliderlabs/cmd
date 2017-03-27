package cli

import (
  "context"
  "sync"
  "encoding/binary"
  "encoding/hex"

  "github.com/spf13/cobra"
)

var defaultCtxRegistry = &ctxRegistry{ctx: make(map[string]context.Context)}

type ctxRegistry struct {
  sync.Mutex
  lastID uint64
  ctx map[string]context.Context
}

func (cr *ctxRegistry) Add(ctx context.Context) string {
  cr.Lock()
  defer cr.Unlock()
  var zero uint64
  cr.lastID = cr.lastID+1 % ^zero
  buf := make([]byte, 4)
  binary.PutUvarint(buf, cr.lastID)
  key := hex.EncodeToString(buf)
  cr.ctx[key] = ctx
  return key
}

func (cr *ctxRegistry) Lookup(key string) context.Context {
  cr.Lock()
  defer cr.Unlock()
  return cr.ctx[key]
}

func (cr *ctxRegistry) Clear(key string) {
  cr.Lock()
  defer cr.Unlock()
  delete(cr.ctx, key)
}

func Context(cmd *cobra.Command) context.Context {
  return defaultCtxRegistry.Lookup(cmd.Root().Annotations["_ctx"])
}
