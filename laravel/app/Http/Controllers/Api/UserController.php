<?php
namespace App\Http\Controllers\Api;
use App\Http\Controllers\GoApiController; // Base controller که sendToGo داره
use Illuminate\Http\Request;
class UserController extends GoApiController
{
    protected $routes = [
        'update' => ['url'=>'/api/user/update','method'=>'PUT'],
        'info'   => ['url'=>'/api/user/info','method'=>'GET'],
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
