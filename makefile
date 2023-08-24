SERVICES = orders-service payments-service

default: lint test security coverage

lint:
	for service in $(SERVICES); do \
		$(MAKE) -C $$service lint; \
	done

test:
	for service in $(SERVICES); do \
		$(MAKE) -C $$service test; \
	done

coverage:
	for service in $(SERVICES); do \
		$(MAKE) -C $$service coverage; \
	done

security:
	for service in $(SERVICES); do \
		$(MAKE) -C $$service security; \
	done

.PHONY: lint test security coverage
