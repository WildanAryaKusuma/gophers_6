import http from 'k6/http';
import { check } from 'k6';
import { Counter } from 'k6/metrics';

// 500 pembeli (50 VU) menyerbu 1 barang lewat HTTP sungguhan.
// Jalankan: k6 run -e ID_BARANG=2 race_test.js
// Hasil yang benar: pembelian_sukses == stok awal, sisanya ditolak (400), 0 error 500.

const pembelianSukses = new Counter('pembelian_sukses');
const ditolakStokHabis = new Counter('ditolak_stok_habis');

export const options = {
  scenarios: {
    serbuan_flash_sale: {
      executor: 'shared-iterations',
      vus: 50,
      iterations: 500,
      maxDuration: '30s',
    },
  },
};

const BASE_URL = 'http://localhost:3000';
const ID_BARANG = Number(__ENV.ID_BARANG || 2);

export default function () {
  const payload = JSON.stringify({ id_barang: ID_BARANG, jumlah: 1, uang: 1000000 });
  const res = http.post(`${BASE_URL}/api/beli`, payload, {
    headers: { 'Content-Type': 'application/json' },
  });

  if (res.status === 200) pembelianSukses.add(1);
  if (res.status === 400) ditolakStokHabis.add(1);

  check(res, {
    'tidak ada error server (500)': (r) => r.status !== 500,
  });
}
