import phonenumbers
from phonenumbers import geocoder, carrier
# Masukkan nomor telepon dalam format internasional
nomor_hp = "+6281252520251"
# Parsing nomor
nomor_parsed = phonenumbers.parse(nomor_hp)
# Menampilkan lokasi/geografi
lokasi = geocoder.description_for_number(nomor_parsed, "id")
print(f"Lokasi: {lokasi}")
# Menampilkan operator seluler
operator = carrier.name_for_number(nomor_parsed, "id")
print(f"Operator: {operator}")
