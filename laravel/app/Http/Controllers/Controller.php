<?php

namespace App\Http\Controllers;
use Illuminate\Routing\Controller;
use Illuminate\Support\Facades\Http;

class GoApiController extends Controller
{
    protected $goBaseUrl;
    public function __construct()
    {
        $this->goBaseUrl = env('GO_API_URL', 'http://go:8080');
    }
    protected function sendToGo(string $endpoint, array $data = [], string $method = 'POST')
    {
        $response = Http::withHeaders([
            'X-API-KEY' => session('api_key') ?? ''
        ])->timeout(10);
        switch(strtoupper($method)) {
            case 'GET':
                $res = $response->get($this->goBaseUrl . $endpoint, $data);
                break;
            case 'POST':
                $res = $response->post($this->goBaseUrl . $endpoint, $data);
            case 'PUT':
                $res = $response->put($this->goBaseUrl . $endpoint, $data);
            case 'DELETE':
                $res = $response->delete($this->goBaseUrl . $endpoint, $data);
            default:
                $res = $response->post($this->goBaseUrl . $endpoint, $data);
        }
        return $res->json();
    }
}