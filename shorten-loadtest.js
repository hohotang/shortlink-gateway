import http from 'k6/http';
import { check, sleep } from 'k6';

export const options = {
    stages: [
        { duration: '10s', target: 100 },
        { duration: '10s', target: 200 },
        { duration: '10s', target: 300 },
        { duration: '10s', target: 400 },
        { duration: '10s', target: 500 },
    ],
};

// 使用__VU作為虛擬用戶ID來確保每個用戶發送的URL是唯一的
export default function () {
    // 每個VU有自己的起始index，基於VU ID進行偏移
    const baseIndex = 1000000;
    const userIndex = baseIndex - (__VU * 1000) - __ITER;

    const url = 'http://localhost:8080/shorten';

    const payload = JSON.stringify({
        original_url: `https://www.ptt.cc/bbs/Gossiping/index${userIndex}.html`,
    });

    const headers = {
        'Content-Type': 'application/json',
        // 'Authorization': 'Bearer YOUR_TOKEN', // Uncomment if authentication is needed
    };

    const res = http.post(url, payload, { headers });

    check(res, {
        'status is 200': (r) => r.status === 200 || r.status === 201,
    });

    sleep(0.1);
}

// █ TOTAL RESULTS

// checks_total.......................: 100093  1997.267423/s
// checks_succeeded...................: 100.00% 100093 out of 100093
// checks_failed......................: 0.00%   0 out of 100093

// ✓ status is 200

// HTTP
// http_req_duration.......................................................: avg=24.36ms min=1.51ms   med=20.05ms  max=197.67ms p(90)=48.5ms   p(95)=61.83ms
//   { expected_response:true }............................................: avg=24.36ms min=1.51ms   med=20.05ms  max=197.67ms p(90)=48.5ms   p(95)=61.83ms
// http_req_failed.........................................................: 0.00%  0 out of 100093
// http_reqs...............................................................: 100093 1997.267423/s

// EXECUTION
// iteration_duration......................................................: avg=124.9ms min=101.82ms med=120.62ms max=298.33ms p(90)=149.08ms p(95)=162.43ms
// iterations..............................................................: 100093 1997.267423/s
// vus.....................................................................: 499    min=10          max=499
// vus_max.................................................................: 500    min=500         max=500

// NETWORK
// data_received...........................................................: 17 MB  344 kB/s
// data_sent...............................................................: 21 MB  421 kB/s