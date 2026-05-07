.PHONY: infra infra-down infra-logs

infra:
	docker compose -f docker-compose.dev.yml up -d

infra-down:
	docker compose -f docker-compose.dev.yml down

infra-logs:
	docker compose -f docker-compose.dev.yml logs -f


