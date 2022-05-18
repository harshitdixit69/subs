package queries

//language=PostgreSQL
const AddCustomerInTransit = `
	INSERT INTO customer_transition ("businessId", "checkoutSessionId", "priceId")
	VALUES (:businessId, :checkoutSessionId, :priceId)
`

//language=PostgreSQL
const AddCustomer = `
	INSERT INTO customers ("stripeCustomerId", "businessId", "subscriptionId", "productId")
	VALUES (:stripeCustomerId, :businessId, :subscriptionId, :productId)
`

//language=PostgreSQL
const AddSubscription = `
	INSERT INTO subscriptions ("stripeSubscriptionId", status, ends_ts)
	VALUES (:stripeSubscriptionId, 'not active', now())
`

//language=PostgreSQL
const UpdateSubscriptionEndDate = `
	UPDATE subscriptions
	SET ends_ts = :end_ts,
		status = 'active'
	WHERE "stripeSubscriptionId" = :stripeSubscriptionId
`

//language=PostgreSQL
const RemoveCustomerInTransition = `
	DELETE FROM customer_transition
	WHERE "checkoutSessionId" = :checkoutSessionId
`

//language=PostgreSQL
const AddInvoice = `
	INSERT INTO invoice_history ("subscriptionId", "customerId", "stripeInvoiceId", paid, total)
	VALUES (:subscriptionId, :customerId, :stripeInvoiceId, :paid, :total)
`

//language=PostgreSQL
const DisableSubscription = `
	UPDATE subscriptions
	SET status='unpaid'
	WHERE "stripeSubscriptionId" = :stripeSubscriptionId
`

//language=PostgreSQL
const SubscriptionChanged = `
	UPDATE customers
	SET "productId" = :productId
	WHERE "subscriptionId" = :subscriptionId
`
