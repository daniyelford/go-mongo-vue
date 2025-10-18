<?php

namespace App\Http\Controllers;

use Illuminate\Routing\Controller;
use Illuminate\Support\Facades\Http;
use Illuminate\Http\Request;

class GoApiController extends Controller
{
    protected string $goBaseUrl;

    public function __construct()
    {
        $this->goBaseUrl = env('GO_API_URL', 'http://go:8080');
    }

    /**
     *
     * @param string $endpoint
     * @param array $data
     * @param string $method
     * @return array
     */
    protected function sendToGo(string $endpoint, array $data = [], string $method = 'POST', Request $request = null): array
    {
        $headers = $request?->headers->all() ?? [];
        $headers = collect($headers)->map(fn($h) => $h[0] ?? '')->toArray();
        $headers['X-API-KEY'] = $headers['X-API-KEY'] ?? '';
        $response = Http::withHeaders($headers)->timeout(10);
        $url = rtrim($this->goBaseUrl, '/') . '/' . ltrim($endpoint, '/');
        $method = strtoupper($method);
        $res = match ($method) {
            'GET' => $response->get($url, $data),
            'PUT' => $response->put($url, $data),
            'DELETE' => $response->delete($url, $data),
            default => $response->post($url, $data),
        };
        return $res->json() ?? ['raw' => $res->body()];
    }
}
