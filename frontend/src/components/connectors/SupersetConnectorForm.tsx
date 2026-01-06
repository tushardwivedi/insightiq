"use client";

import { useState, memo, useCallback, useMemo } from "react";
import { DataConnector, SupersetConnectorConfig } from "@/types";
import { apiClient } from "@/lib/api";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Textarea } from "@/components/ui/textarea";
import { Label } from "@/components/ui/label";
import {
  Card,
  CardContent,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card";
import { RadioGroup, RadioGroupItem } from "@/components/ui/radio-group";
import { Alert, AlertDescription } from "@/components/ui/alert";
import { Loader2, CheckCircle, XCircle } from "lucide-react";

interface SupersetConnectorFormProps {
  onCancel: () => void;
  onSuccess: (connector: DataConnector) => void;
  connector?: DataConnector;
}

function SupersetConnectorForm({
  onCancel,
  onSuccess,
  connector,
}: SupersetConnectorFormProps) {
  const [formData, setFormData] = useState({
    name: connector?.name || "",
    url: connector?.config.url || "",
    username: (connector?.config as SupersetConnectorConfig)?.username || "",
    password: (connector?.config as SupersetConnectorConfig)?.password || "",
    bearer_token:
      (connector?.config as SupersetConnectorConfig)?.bearer_token || "",
  });
  const [loading, setLoading] = useState(false);
  const [testing, setTesting] = useState(false);
  const [testResult, setTestResult] = useState<{
    success: boolean;
    message: string;
  } | null>(null);
  const [authMethod, setAuthMethod] = useState<"credentials" | "token">(
    "credentials"
  );

  // Memoized change handlers
  const handleNameChange = useCallback(
    (e: React.ChangeEvent<HTMLInputElement>) => {
      setFormData((prev) => ({ ...prev, name: e.target.value }));
    },
    []
  );

  const handleUrlChange = useCallback(
    (e: React.ChangeEvent<HTMLInputElement>) => {
      setFormData((prev) => ({ ...prev, url: e.target.value }));
    },
    []
  );

  const handleUsernameChange = useCallback(
    (e: React.ChangeEvent<HTMLInputElement>) => {
      setFormData((prev) => ({ ...prev, username: e.target.value }));
    },
    []
  );

  const handlePasswordChange = useCallback(
    (e: React.ChangeEvent<HTMLInputElement>) => {
      setFormData((prev) => ({ ...prev, password: e.target.value }));
    },
    []
  );

  const handleBearerTokenChange = useCallback(
    (e: React.ChangeEvent<HTMLTextAreaElement>) => {
      setFormData((prev) => ({ ...prev, bearer_token: e.target.value }));
    },
    []
  );

  const handleAuthMethodChange = useCallback(
    (value: "credentials" | "token") => {
      setAuthMethod(value);
    },
    []
  );

  const handleCredentialsRadioChange = useCallback(() => {
    handleAuthMethodChange("credentials");
  }, [handleAuthMethodChange]);

  const handleTokenRadioChange = useCallback(() => {
    handleAuthMethodChange("token");
  }, [handleAuthMethodChange]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setTestResult(null);

    try {
      const config: SupersetConnectorConfig = {
        url: formData.url.trim(),
        username: authMethod === "credentials" ? formData.username.trim() : "",
        password: authMethod === "credentials" ? formData.password : "",
        bearer_token:
          authMethod === "token" ? formData.bearer_token.trim() : "",
      };

      const connectorData = {
        name: formData.name.trim(),
        type: "superset" as const,
        config,
      };

      let result: DataConnector;
      if (connector) {
        result = await apiClient.updateConnector(connector.id, connectorData);
      } else {
        result = await apiClient.createConnector(connectorData);
      }

      onSuccess(result);
    } catch (error: any) {
      console.error("Failed to save connector:", error);
      setTestResult({
        success: false,
        message: error.message || "Failed to save connector. Please try again.",
      });
    } finally {
      setLoading(false);
    }
  };

  const handleTestConnection = async () => {
    setTesting(true);
    setTestResult(null);

    try {
      const config: SupersetConnectorConfig = {
        url: formData.url.trim(),
        username: authMethod === "credentials" ? formData.username.trim() : "",
        password: authMethod === "credentials" ? formData.password : "",
        bearer_token:
          authMethod === "token" ? formData.bearer_token.trim() : "",
      };

      const result = await apiClient.testConnectorConfig({
        type: "superset",
        config,
      });

      setTestResult({
        success: result.success,
        message:
          result.message ||
          (result.success ? "Connection successful!" : "Connection failed"),
      });
    } catch (error: any) {
      console.error("Connection test failed:", error);
      setTestResult({
        success: false,
        message:
          error.message ||
          "Connection test failed. Please check your settings.",
      });
    } finally {
      setTesting(false);
    }
  };

  const isFormValid = useMemo(() => {
    return (
      formData.name.trim() &&
      formData.url.trim() &&
      (authMethod === "credentials"
        ? formData.username.trim() && formData.password.trim()
        : formData.bearer_token.trim())
    );
  }, [
    formData.name,
    formData.url,
    formData.username,
    formData.password,
    formData.bearer_token,
    authMethod,
  ]);

  return (
    <Card>
      <CardHeader>
        <CardTitle className="flex items-center gap-2">
          📊 {connector ? "Edit" : "Add"} Superset Connector
        </CardTitle>
        <CardDescription>
          Connect to your Apache Superset instance for analytics data
        </CardDescription>
      </CardHeader>

      <CardContent>
        <form onSubmit={handleSubmit} className="space-y-6">
          {/* Name */}
          <div className="space-y-2">
            <Label htmlFor="name">Connection Name *</Label>
            <Input
              type="text"
              id="name"
              value={formData.name}
              onChange={handleNameChange}
              placeholder="My Superset Instance"
              required
            />
          </div>

          {/* URL */}
          <div className="space-y-2">
            <Label htmlFor="url">Superset URL *</Label>
            <Input
              type="url"
              id="url"
              value={formData.url}
              onChange={handleUrlChange}
              placeholder="https://superset.example.com"
              required
            />
          </div>

          {/* Authentication Method */}
          <div className="space-y-3">
            <Label>Authentication Method</Label>
            <RadioGroup
              value={authMethod}
              onValueChange={handleAuthMethodChange}
            >
              <div className="flex items-center space-x-2">
                <RadioGroupItem value="credentials" id="credentials" />
                <Label htmlFor="credentials">Username & Password</Label>
              </div>
              <div className="flex items-center space-x-2">
                <RadioGroupItem value="token" id="token" />
                <Label htmlFor="token">Bearer Token</Label>
              </div>
            </RadioGroup>
          </div>

          {/* Authentication Fields */}
          {authMethod === "credentials" ? (
            <>
              <div className="space-y-2">
                <Label htmlFor="username">Username *</Label>
                <Input
                  type="text"
                  id="username"
                  value={formData.username}
                  onChange={handleUsernameChange}
                  placeholder="admin"
                  required={authMethod === "credentials"}
                />
              </div>

              <div className="space-y-2">
                <Label htmlFor="password">Password *</Label>
                <Input
                  type="password"
                  id="password"
                  value={formData.password}
                  onChange={handlePasswordChange}
                  placeholder="••••••••"
                  required={authMethod === "credentials"}
                />
              </div>
            </>
          ) : (
            <div className="space-y-2">
              <Label htmlFor="bearer_token">Bearer Token *</Label>
              <Textarea
                id="bearer_token"
                value={formData.bearer_token}
                onChange={handleBearerTokenChange}
                placeholder="eyJhbGciOiJIUzI1NiIsInR5cCI6..."
                rows={3}
                className="resize-none"
                required={authMethod === "token"}
              />
              <p className="text-xs text-muted-foreground">
                Generate a bearer token from Superset's API settings
              </p>
            </div>
          )}

          {/* Test Result */}
          {testResult && (
            <Alert variant={testResult.success ? "default" : "destructive"}>
              <div className="flex items-center gap-2">
                {testResult.success ? (
                  <CheckCircle className="h-4 w-4" />
                ) : (
                  <XCircle className="h-4 w-4" />
                )}
                <AlertDescription>{testResult.message}</AlertDescription>
              </div>
            </Alert>
          )}

          {/* Actions */}
          <div className="flex flex-wrap gap-3 pt-4">
            <Button
              type="button"
              onClick={handleTestConnection}
              disabled={!isFormValid || testing}
              variant="outline"
              className="flex-1"
            >
              {testing && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
              {testing ? "Testing..." : "Test Connection"}
            </Button>
            <br />
            <Button type="button" onClick={onCancel} variant="outline">
              Cancel
            </Button>
            <Button type="submit" disabled={!isFormValid || loading}>
              {loading && <Loader2 className="mr-2 h-4 w-4 animate-spin" />}
              {loading ? "Saving..." : connector ? "Update" : "Save"}
            </Button>
          </div>
        </form>
      </CardContent>
    </Card>
  );
}

export default memo(SupersetConnectorForm);
