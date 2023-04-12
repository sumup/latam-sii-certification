# latam-sii-certification

The purpose of this piece of code is to serve as a tool for sii's certification process, here you can realize the whole certification process by just running this program.

## Certification cases

| Document Type | Quantity | Amount           |
|---------------|----------|------------------|
| 48 (Voucher)  | 4        | Distinct to Zero |
| 48 (Voucher)  | 3        | Zero             |
| 33 (Bill)     | 1        | Distinct to Zero |
| 00 (Unkown)   | 1        | Distinct to Zero |
| 99 (No sell)  | 1        | Distinct to Zero |

Considerations: 
- External Track ID must be zero
- Number of transactions must be greater than zero
- Channel:
    - 0 for CP
    - 1 for CNP
- VatId (RUT) must be valid 

## How to run 

make run