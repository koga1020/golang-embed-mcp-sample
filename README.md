# Go Embed MCP サンプル

Go の `//go:embed` ディレクティブを使用して、静的アセットを埋め込んだ MCP (Model Context Protocol) サーバーの構築方法を示すサンプルプロジェクト。

## このサンプルで実演すること

- **静的埋め込み**: プロンプトとリソースをビルド時にバイナリに直接埋め込み
- **ゼロ依存**: デプロイメント時に外部ファイルが不要  
- **MCP統合**: 標準MCPプロトコルを通じて埋め込みコンテンツを配信
- **選択的読み込み**: コマンドラインフラグによるコンテンツフィルタリング

## 主な特徴

### Embed ディレクティブの使用
```go
//go:embed prompts/*
var embeddedPrompts embed.FS

//go:embed resources/*  
var embeddedResources embed.FS
```

## ビルド & テスト

```bash
# サンプルをビルド
go build -o embed-mcp cmd/embed-mcp/main.go

# 全コンテンツでテスト
./embed-mcp

# フィルタリングしてテスト  
./embed-mcp --prompts demo --resources config
```

## プロジェクト構造

```
cmd/embed-mcp/
├── prompts/          # プロンプト
│   ├── demo.md
│   └── sample.md
└── resources/        # リソース
    ├── config.json
    ├── schema.yaml
    └── info.md
```
