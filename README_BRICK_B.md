# Brick B — Paid Listing Toggle (private/public + fee trigger model)

## Scope
Listing visibility/payment is NOT part of asset identity.
No payment processing. No search. No ranking. No authority claims.

## Endpoint
POST /set_listing

### Request
{
  "asset_id": "asset:<sha256>",
  "visibility": "private" | "public",
  "paid": true | false
}

### Response
{
  "asset_id": "...",
  "listing_status": "private" | "payment_required" | "public",
  "fee_required": true | false,
  "created_at_utc": "<server time>"
}

## Fee trigger model (minimal)
- visibility="private" => listing_status="private", fee_required=false
- visibility="public"  => fee_required=true
  - paid=false => listing_status="payment_required"
  - paid=true  => listing_status="public"

## Invariants proven
- Listing state does not change or recompute asset_id
- Public listings require fee (fee_required=true) even if unpaid
- Private listings are free (fee_required=false)
