CREATE TABLE users (
  id UUID PRIMARY KEY,
  handle TEXT NOT NULL UNIQUE,
  display_name TEXT NOT NULL,
  email TEXT NOT NULL UNIQUE,
  password_hash TEXT NOT NULL,
  rank TEXT NOT NULL DEFAULT 'Outer Disciple',
  xp INTEGER NOT NULL DEFAULT 0,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE problems (
  id UUID PRIMARY KEY,
  slug TEXT NOT NULL UNIQUE,
  title TEXT NOT NULL,
  difficulty TEXT NOT NULL,
  gate TEXT NOT NULL,
  summary TEXT NOT NULL,
  body TEXT NOT NULL,
  constraints JSONB NOT NULL DEFAULT '[]',
  examples JSONB NOT NULL DEFAULT '[]',
  tags TEXT[] NOT NULL DEFAULT '{}',
  acceptance_rate INTEGER NOT NULL DEFAULT 0,
  xp INTEGER NOT NULL DEFAULT 25,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE starter_code (
  id UUID PRIMARY KEY,
  problem_id UUID NOT NULL REFERENCES problems(id) ON DELETE CASCADE,
  language_id INTEGER NOT NULL,
  language TEXT NOT NULL,
  code TEXT NOT NULL
);

CREATE TABLE test_cases (
  id UUID PRIMARY KEY,
  problem_id UUID NOT NULL REFERENCES problems(id) ON DELETE CASCADE,
  name TEXT NOT NULL,
  stdin TEXT NOT NULL,
  expected_output TEXT NOT NULL,
  hidden BOOLEAN NOT NULL DEFAULT true
);

CREATE TABLE submissions (
  id UUID PRIMARY KEY,
  user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
  problem_id UUID NOT NULL REFERENCES problems(id) ON DELETE CASCADE,
  language_id INTEGER NOT NULL,
  language TEXT NOT NULL,
  source_code TEXT NOT NULL,
  status TEXT NOT NULL,
  score INTEGER NOT NULL DEFAULT 0,
  runtime_ms INTEGER NOT NULL DEFAULT 0,
  memory_kb INTEGER NOT NULL DEFAULT 0,
  created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE submission_results (
  id UUID PRIMARY KEY,
  submission_id UUID NOT NULL REFERENCES submissions(id) ON DELETE CASCADE,
  test_name TEXT NOT NULL,
  status TEXT NOT NULL,
  stdout TEXT,
  stderr TEXT,
  compile_output TEXT,
  runtime_ms INTEGER NOT NULL DEFAULT 0,
  memory_kb INTEGER NOT NULL DEFAULT 0
);
