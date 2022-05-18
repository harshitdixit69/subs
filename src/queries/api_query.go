package queries

// language=PostgreSQL
const MatchingKey = `
	SELECT ( CASE WHEN api_key.key = :apiKey THEN TRUE ELSE false END) as auth
	FROM api_key
`
