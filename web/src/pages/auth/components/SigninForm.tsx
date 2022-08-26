import type { HTMLChakraProps } from "@chakra-ui/react";
import {
  Button,
  Center,
  chakra,
  FormControl,
  FormLabel,
  Input,
  Stack,
} from "@chakra-ui/react";
import { useAtom } from "jotai";
import { useState } from "react";
import { useLocation, useNavigate } from "react-router-dom";

import AuthApi from "../../../api/auth-api";
import type { ShortUser } from "../../../models/user";
import { jwtToken, refreshJwtToken } from "../../../store";
import { useErrorHandler } from "../../../utils/handle-get-error";

import Link from "./Link";
import PasswordField from "./PasswordField";

const SigninForm = (props: HTMLChakraProps<"form">) => {
  const [, setJwt] = useAtom(jwtToken);
  const [, setRefreshJwt] = useAtom(refreshJwtToken);
  const [user, setUser] = useState<ShortUser>({
    email: "",
    password: "",
  });
  const errorHandler = useErrorHandler();
  const navigate = useNavigate();
  const location = useLocation();
  // @ts-ignore
  const from = location.state?.from || "/smers";

  return (
    <chakra.form
      onSubmit={(e) => {
        e.preventDefault();
        AuthApi.signin(user)
          .then((payload) => {
            setJwt(payload.token);
            setRefreshJwt(payload.refreshToken);
            navigate(from, { replace: true });
          })
          .catch(errorHandler);
      }}
      {...props}
    >
      <Stack spacing="6">
        <FormControl id="email">
          <FormLabel>Email address</FormLabel>
          <Input
            value={user.email}
            onChange={(event) =>
              setUser({ ...user, email: event.target.value })
            }
            name="email"
            type="email"
            autoComplete="email"
            required
          />
        </FormControl>
        <PasswordField
          value={user.password}
          onChange={(event) =>
            setUser({ ...user, password: event.target.value })
          }
        />
        <Center>
          <Link href="/reset-password">
            Forgot password?
          </Link>
        </Center>
        <Button type="submit" colorScheme="blue" size="lg" fontSize="md">
          Sign in
        </Button>
      </Stack>
    </chakra.form>
  );
};
export default SigninForm;
