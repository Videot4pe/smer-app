import { useAtom } from "jotai";
import { Navigate, useLocation } from "react-router-dom";

import { isLoggedIn } from "../../store";

type PrivateRouteProps = {
  children: React.ReactNode;
  redirectTo?: string;
};

const RequireAuth = ({
  children,
  redirectTo = "/signin",
}: PrivateRouteProps) => {
  const [isAuthenticated] = useAtom(isLoggedIn);
  const location = useLocation();

  return isAuthenticated ? (
    (children as React.ReactElement)
  ) : (
    <Navigate to={redirectTo} state={{ from: location }} replace />
  );
};

export default RequireAuth;
