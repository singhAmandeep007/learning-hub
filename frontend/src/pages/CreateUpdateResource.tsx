import { CreateUpdateResourceForm } from "./components/CreateUpdateResourceForm";

export function CreateUpdateResource() {
  const onSuccess = () => {
    console.log("Resource saved:");
  };

  const onCancel = () => {
    console.log("Cancelled");
  };

  return (
    <div>
      <CreateUpdateResourceForm
        onCancel={onCancel}
        onSuccess={onSuccess}
      />
    </div>
  );
}
