.PHONY: help install-backend install-frontend build-backend build-frontend test-backend test-frontend run-backend run-frontend clean

# デフォルトターゲット
help:
	@echo "AWS S3 Knowledge Chatbot - Available targets:"
	@echo ""
	@echo "Infrastructure:"
	@echo "  terraform-init    - Terraformを初期化"
	@echo "  terraform-plan    - Terraformプランを表示"
	@echo "  terraform-apply   - インフラをデプロイ"
	@echo "  terraform-destroy - インフラを削除"
	@echo ""
	@echo "Backend (Go):"
	@echo "  install-backend   - バックエンドの依存関係をインストール"
	@echo "  build-backend     - バックエンドをビルド"
	@echo "  test-backend      - バックエンドのテストを実行"
	@echo "  run-backend       - バックエンドを起動"
	@echo ""
	@echo "Frontend (Angular):"
	@echo "  install-frontend  - フロントエンドの依存関係をインストール"
	@echo "  build-frontend    - フロントエンドをビルド"
	@echo "  test-frontend     - フロントエンドのテストを実行"
	@echo "  run-frontend      - フロントエンドを起動 (開発サーバー)"
	@echo ""
	@echo "All:"
	@echo "  install           - すべての依存関係をインストール"
	@echo "  build             - すべてをビルド"
	@echo "  test              - すべてのテストを実行"
	@echo "  clean             - ビルド成果物を削除"

# Infrastructure
terraform-init:
	cd infrastructure && terraform init

terraform-plan:
	cd infrastructure && terraform plan

terraform-apply:
	cd infrastructure && terraform apply

terraform-destroy:
	cd infrastructure && terraform destroy

# Backend
install-backend:
	cd backend && go mod download

build-backend:
	cd backend && go build -o server main.go

test-backend:
	cd backend && go test -v ./...

run-backend:
	cd backend && go run main.go

# Frontend
install-frontend:
	cd frontend && npm install

build-frontend:
	cd frontend && npm run build

test-frontend:
	cd frontend && npm test

run-frontend:
	cd frontend && npm start

# All
install: install-backend install-frontend

build: build-backend build-frontend

test: test-backend test-frontend

# Clean
clean:
	rm -f backend/server
	rm -rf frontend/dist
	rm -rf frontend/.angular
	@echo "Clean complete"
