给我outputs时，要确认这些没有被花费过。（即避免同一outputs给我多次）

/wallet/committransactions:
[
  transactions....
]
json response:
{
  "transactionids": [
    "1234567890abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
    "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
    "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"
  ]
}

/wallet/checkoutput
{"transaction": ....}
json response:
{
  "spendable":   [1, 2],
  "unspendable": [0]
}
