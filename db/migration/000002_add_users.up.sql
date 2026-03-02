-- 添加 users 表
CREATE TABLE "users" (
  "username" varchar PRIMARY KEY,
  "hashed_password" varchar NOT NULL,
  "full_name" varchar NOT NULL,
  "email" varchar UNIQUE NOT NULL,
  "password_changed_at" timestamptz NOT NULL DEFAULT '0001-01-01 00:00:00Z',
  "created_at" timestamptz NOT NULL DEFAULT (now())
);
-- 添加外键
ALTER TABLE "accounts" ADD FOREIGN KEY ("owner") REFERENCES "users" ("username") DEFERRABLE INITIALLY IMMEDIATE;

-- 添加 unique constraints 唯一约束
-- CREATE UNIQUE INDEX ON "accounts" ("owner", "currency");
ALTER TABLE "accounts" ADD CONSTRAINT "owner_currency_key" UNIQUE ("owner", "currency");

-- add users table with unique & foreign key constraints in PostgresSQL
-- 把 dbdiagram 导出的 PostgresSQL 代码把补充的内容粘贴进来
-- 全部覆盖是不好的
-- 这个文件来源于在 terminal 运行 migrate create -ext sql -dir db/migration -seq add_users