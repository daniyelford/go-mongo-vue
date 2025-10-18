<?php

namespace App\Http\Controllers\Api;

use App\Http\Controllers\GoApiController;
use Illuminate\Http\Request;

class AuthController extends GoApiController
{
    protected $routes = [
        'hasFingerPrint' => ['url'=>'/api/login/fingerPrint/has','method'=>'POST'],
        'sendCode' => ['url'=>'/api/login/send','method'=>'POST'],
        'verifyCode' => ['url'=>'/api/login/verify','method'=>'POST'],
        'loginFingerPrintStart' => ['url'=>'/api/login/fingerPrint/start','method'=>'POST'],
        'loginFingerPrintEnd' => ['url'=>'/api/login/fingerPrint/end','method'=>'POST'],
        'refreshToken' => ['url'=>'/api/token/refresh','method'=>'POST'],
        'logout' => ['url'=>'/api/auth/logout','method'=>'GET'],
        'validateToken' => ['url'=>'/api/auth/validate','method'=>'GET'],
        'register' => ['url'=>'/api/register/save','method'=>'POST'],
        'registerFingerPrintStart' => ['url'=>'/api/register/fingerPrint/start','method'=>'POST'],
        'registerFingerPrintEnd' => ['url'=>'/api/register/fingerPrint/end','method'=>'POST'],
    ];
    public function call(Request $request, $action)
    {
        if (!isset($this->routes[$action])) {
            return response()->json(['status' => 'error', 'message' => 'متد نامعتبر']);
        }
        $endpoint = $this->routes[$action]['url'];
        $method = $this->routes[$action]['method'];
        return $this->sendToGo($endpoint, $request->all(),$method,$request);
    }
}
