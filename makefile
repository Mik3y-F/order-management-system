SERVICES = orders-service

default: lint test security

lint:
	for service in $(SERVICES); do \
		$(MAKE) -C $$service lint; \
	done

test:
	for service in $(SERVICES); do \
		$(MAKE) -C $$service test; \
	done

security:
	for service in $(SERVICES); do \
		$(MAKE) -C $$service security; \
	done

.PHONY: lint test security
