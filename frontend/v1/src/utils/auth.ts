import { backendRoutes, ContentType, methods, statusCodes } from "@/constants";
import { AuthStatus } from "@/models/state";

export async function isAuthenticated(): Promise<AuthStatus> {
  const status = await fetch(backendRoutes.USER_AUTH_STATUS, {
    method: methods.GET,
    headers: {
      "Content-Type": ContentType,
    },
  }).then((res: Response) => {
    if (
      res.status === statusCodes.OK ||
      res.status === statusCodes.UNAUTHORIZED
    ) {
      return res.json();
    } else {
      console.log(`invalid response code: ${res.status}`);
      console.dir(res);
    }
  });

  if (status) {
    return {
      isAuthenticated: status["isAuthenticated"],
      isAdmin: status["isAdmin"],
    } as AuthStatus;
  }

  throw "unauthorized";
}
