import { useToast } from "@chakra-ui/react";

export const useSuccessHandler = (title: string) => {
  const toast = useToast();

  return (message: string) => {
    toast({
      title,
      description: message,
      status: "success",
      duration: 2000,
      isClosable: true,
    });
  };
};
