import {
    Box,
    Heading,
    Text,
    useToast,
  } from "@chakra-ui/react";
  
  import Card from "./components/Card";
  import Logo from "./components/Logo";
  import ChangePasswordForm from "./components/ChangePasswordForm";
import { useNavigate, useSearchParams } from "react-router-dom";
  
  const ChangePassword = () => {
    const toast = useToast();
    const navigate = useNavigate();
    const [searchParams] = useSearchParams();
    const token = searchParams.get("token")

    if (!token) {
      navigate('/signin', { replace: true });
      toast({
        title: 'Empty token',
        description: 'Please, try to reset email again',
        status: "error",
        isClosable: true,
      });

      return
    }

    return (<Box bg="pink.100" minH="100vh" py="12" px={{ base: "4", lg: "8" }}>
      <Box maxW="md" mx="auto">
        <Logo mx="auto" h="8" mb={{ base: "10", md: "20" }} />
        <Heading textAlign="center" size="xl" fontWeight="extrabold">
          Reset password
        </Heading>
        <Text mt="4" mb="8" align="center" maxW="md" fontWeight="medium">
          <Text as="span">Set your new password</Text>
        </Text>
        <Card>
          <ChangePasswordForm token={ token }/>
        </Card>
      </Box>
    </Box>
  )};
  
  export default ChangePassword;
