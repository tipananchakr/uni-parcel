  "use client";

  import { useColorMode, useColorModeValue } from "@/components/ui/color-mode";
  import { Button, Field, Input } from "@chakra-ui/react"
  import { Moon, Sun } from "lucide-react";
  import { useForm } from "react-hook-form";

  interface LoginInterface {
    Email: string
    Password: string
  }

  export default function Home() {
    const { register, handleSubmit, formState: { errors, isValid }, reset } = useForm<LoginInterface>();

    const { toggleColorMode } = useColorMode();

    const onSubmit = (data: LoginInterface) => {
      console.log(data);
    }

    return (
      <div>

        <div className="border h-screen flex items-center justify-center">
          <Button cursor={"pointer"} onClick={toggleColorMode}>{useColorModeValue(<Moon />, <Sun />)}</Button>

          <div>
            <h1 className="text-2xl font-bold">Login</h1>
            <form onSubmit={handleSubmit(onSubmit)}>
              <Field.Root invalid={!!errors.Email && !!errors.Password}>
                <Input
                  type="email"
                  placeholder="email"
                  {...register("Email", {
                    required: "Email is required"
                  })}
                />
                <Field.ErrorText>{errors.Email?.message}</Field.ErrorText>
                <Input
                  type="password"
                  placeholder="password"
                  {...register("Password", {
                    required: "Password is required",
                    minLength: {
                      value: 6,
                      message: "Password must be at least 6 characters"
                    }
                  })}
                />
                <Field.ErrorText>{errors.Password?.message}</Field.ErrorText>
              </Field.Root>


              <Button type="submit">Login</Button>
            </form>
          </div>

        </div>
      </div>
    );
  }
