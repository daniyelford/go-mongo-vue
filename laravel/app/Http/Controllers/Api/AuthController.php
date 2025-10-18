<?php

namespace App\Http\Controllers\Api;

use App\Http\Controllers\GoApiController;
use Illuminate\Http\Request;

class AuthController extends GoApiController
{
    protected $routes = [
        'hasFingerPrint' => '/api/login/fingerPrint/has',
        'sendCode' => '/api/login/send',
        'verifyCode' => '/api/login/verify',
        'loginFingerPrintStart' => '/api/login/fingerPrint/start',
        'loginFingerPrintEnd' => '/api/login/fingerPrint/end',
        'refreshToken' => '/api/token/refresh',
        'logout' => '/api/auth/logout',
        'validateToken' => '/api/auth/validate',
        'register' => '/api/register/save',
        'registerFingerPrintStart' => '/api/register/fingerPrint/start',
        'registerFingerPrintEnd' => '/api/register/fingerPrint/end',
    ];
    public function call(Request $request, $action)
    {
        if (!isset($this->routes[$action])) {
            return response()->json(['status' => 'error', 'message' => 'متد نامعتبر']);
        }
        $endpoint = $this->routes[$action];
        return $this->sendToGo($endpoint, $request->all());
    }
}
