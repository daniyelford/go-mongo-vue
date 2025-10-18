<?php

namespace App\Http\Controllers\Api;

use App\Http\Controllers\GoApiController; // Base controller که sendToGo داره
use Illuminate\Http\Request;

class AuthController extends GoApiController
{
    // آرایه endpoint ها: key = نام متد، value = endpoint Go
    protected $routes = [
        'hasFingerPrint' => '/api/auth/has-finger-print',
        'sendCode' => '/api/auth/send-code',
        'verifyCode' => '/api/auth/verify-code',
        'loginFingerPrintStart' => '/api/auth/login-finger-start',
        'loginFingerPrintEnd' => '/api/auth/login-finger-end',
        'refreshToken' => '/api/auth/refresh-token',
        'logout' => '/api/auth/logout',
        'validateToken' => '/api/auth/validate-token',
        'register' => '/api/auth/register',
        'registerFingerPrintStart' => '/api/auth/register-finger-start',
        'registerFingerPrintEnd' => '/api/auth/register-finger-end',
    ];

    // فانکشن عمومی که همه درخواست‌ها ازش عبور می‌کنن
    public function call(Request $request, $action)
    {
    return response()->json(['status' => 'OK'], 200);

        if (!isset($this->routes[$action])) {
            return response()->json(['status' => 'error', 'message' => 'متد نامعتبر']);
        }

        $endpoint = $this->routes[$action];
        return $this->sendToGo($endpoint, $request->all());
    }
}
