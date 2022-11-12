START TRANSACTION
;

CREATE TABLE IF NOT EXISTS account(
	id VARCHAR(256) PRIMARY KEY,
	parent_id VARCHAR(36),
	name VARCHAR(255) UNIQUE,
	type VARCHAR(64),
	basis VARCHAR(6),
	created_at TIMESTAMP DEFAULT NOW()
)
;

CREATE TABLE IF NOT EXISTS transaction(
	id VARCHAR(256) PRIMARY KEY,
	description TEXT,
	created_at TIMESTAMP DEFAULT NOW()
)
;

CREATE TABLE IF NOT EXISTS entry(
	id VARCHAR(256) PRIMARY KEY,
	transaction_id VARCHAR(256),
	description TEXT,
	debit_account VARCHAR(256),
	credit_account VARCHAR(256),
	amount_value BIGINT,
	amount_currency CHAR(3),
	created_at TIMESTAMP DEFAULT NOW(),
	CONSTRAINT fk_transaction FOREIGN KEY(transaction_id) REFERENCES transaction(id),
	CONSTRAINT fk_debit_account FOREIGN KEY(debit_account) REFERENCES account(id),
	CONSTRAINT fk_credit_account FOREIGN KEY(credit_account) REFERENCES account(id)
)
;

COMMIT
;