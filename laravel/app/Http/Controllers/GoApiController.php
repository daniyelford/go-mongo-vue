<?php

namespace App\Http\Controllers;

use Illuminate\Routing\Controller;
use Illuminate\Support\Facades\Http;

class GoApiController extends Controller
{
    protected string $goBaseUrl;

    public function __construct()
    {
        $this->goBaseUrl = env('GO_API_URL', 'http://go:8080');
    }

    /**
     * ارسال درخواست به API Go
     *
     * @param string $endpoint
     * @param array $data
     * @param string $method
     * @return array
     */
    protected function sendToGo(string $endpoint, array $data = [], string $method = 'POST'): array
    {
        $response = Http::withHeaders([
            'X-API-KEY' => session('api_key') ?? ''
        ])->timeout(10);

        switch (strtoupper($method)) {
            case 'GET':
                $res = $response->get(rtrim($this->goBaseUrl, '/') . '/' . ltrim($endpoint, '/'), $data);
                break;

            case 'POST':
            default:
                $res = $response->post(rtrim($this->goBaseUrl, '/') . '/' . ltrim($endpoint, '/'), $data);
        }

        return $res->json() ?? [];
    }
}
