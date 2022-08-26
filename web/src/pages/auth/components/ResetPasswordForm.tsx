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
import { useState } from "react";

import AuthApi from "../../../api/auth-api";
import { useSuccessHandler } from "../../../utils/handle-success";
import { useErrorHandler } from "../../../utils/handle-get-error";

const ResetPasswordForm = (props: HTMLChakraProps<"form">) => {
  const [email, setEmail] = useState<string>('');
  const errorHandler = useErrorHandler();
  const successHandler = useSuccessHandler('Link for changing password was sent');

  return (
    <chakra.form
      onSubmit={(e) => {
        e.preventDefault();
        AuthApi.reset(email)
          .then(() => successHandler('Check your email address'))
          .catch(errorHandler);
      }}
      {...props}
    >
      <Stack spacing="6">
        <FormControl id="email">
          <FormLabel>Email address</FormLabel>
          <Input
            value={email}
            onChange={(event) =>
              setEmail(event.target.value)
            }
            name="email"
            type="email"
            autoComplete="email"
            required
          />
        </FormControl>
        <Button type="submit" colorScheme="blue" size="lg" fontSize="md">
          Reset password
        </Button>
      </Stack>
    </chakra.form>
  );
};
export default ResetPasswordForm;
