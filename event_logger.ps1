# check for new events every hour, POST request them to server
while ($true) {
  # Check new events in the past X amount of time
  Write-Host "Running task at: $(Get-Date)"
  $events = Get-WinEvent -LogName Application | Sort-Object TimeCreated -Descending | Select-Object -First 5

  $events | Format-Table -Property TimeCreated, Id, LevelDisplayName, Message -AutoSize

  $url = "http://localhost:3333/logEvent"
  $data = $events | ForEach-Object {
    @{
      TimeCreated = $_.TimeCreated
      Id = $_.Id
      Severity = $_.LevelDisplayName
      Message = $_.Message
    }
  } | ConvertTo-Json -Depth 10


  Invoke-RestMethod -Uri $url -Method Post -Headers @{ "Content-Type" = "application/json" } -Body $data

  # Sleep for 1 hour (900 seconds)
  Start-Sleep -Seconds 900 
}