package queries

//language=PostgreSQL
const GetSubscriptionId = `
	select id as "subscriptionId" from subscriptions
	where "stripeSubscriptionId" = :stripeSubscriptionId
`

//language=PostgreSQL
const GetBusinessAndPriceId = `
	SELECT "businessId", "priceId" from customer_transition
	WHERE "checkoutSessionId" = :checkoutSessionId
`

// language=PostgreSQL
const GetCustomerId = `
	SELECT id as "customerId" FROM customers
	WHERE "stripeCustomerId" = :stripeCustomerId
`

//language=PostgreSQL
const GetStripeCustomerId = `
	SELECT "stripeCustomerId" FROM customers
	WHERE "businessId" = :businessId
`

//language=PostgreSQL
const GetProductId = `
	SELECT id AS "productId"
	FROM products
	WHERE "priceId" = :priceId
`

//language=PostgreSQL
const GetFreeTrialCode = `
	SELECT id as "freeTrialCodeId"
	FROM free_trials
	WHERE code = :freeTrialCode
`

//language=PostgreSQL
const CheckFreeTrialCode = `
	SELECT EXISTS(
		SELECT id 
		FROM free_trials
		WHERE code = :freeTrialCode
	)
`
