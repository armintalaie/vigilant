<script lang="ts">
  import {GetLogs} from '../wailsjs/go/main/App.js'

  let name: string
  let logs: any[] = []
  let columns = [
    {
      id: "id",
      type: "number",
      name: "ID",
      minWidth: 100
    },
    {
      id: "message",
      type: "string",
      name: "Message",
      minWidth: 300
    },
    {
      id: "level",
      type: "string",
      name: "Level",
      minWidth: 100
    },
    {
      id: "severity",
      type: "number",
      name: "Severity",
      minWidth: 100
    },
    {
      id: "timestamp",
      type: "number",
      name: "Timestamp",
      minWidth: 180
    },
    {
      id: "origin",
      type: "string",
      name: "Origin",
      minWidth: 150
    },
    {
      id: "source",
      type: "string",
      name: "Source",
      minWidth: 160
    },
    {
      id: "type",
      type: "string",
      name: "Type",
      minWidth: 100
    },
    {
      id: "group",
      type: "string",
      name: "Group",
      minWidth: 100
    },
    {
      id: "tags",
      type: "string",
      name: "Tags",
      minWidth: 150
    },
    {
      id: "data",
      type: "object",
      name: "Data",
      minWidth: 200
    }
  ]

    async function GetLogsRequest() {
        const res = await GetLogs()
        logs = res == null ? [] : res
    }

  function greet(): void {
    // Greet(name).then(result => resultText = result)
  }


</script>

<main>
  <div class="control">
    <div class="options">
    <input type="text" bind:value={name} placeholder="Search" />
        <button on:click={greet}>Query</button>
    </div>
    <div class="options">
      <button on:click={GetLogsRequest}>Refresh</button>
    </div>
  </div>



  <section>
        <table>
          <thead>
          <tr>
            {#each columns as column}
              <th
                      class="text-left p-2 border-b"
                      style:min-width="{column.minWidth}px"
              >
                {column.name}
              </th>
            {/each}
          </tr>
          </thead>
          <tbody>
          {#each logs as log}
            <tr
                    class="hover:bg-gray-100"
                    style:background-color={
            log.level === 2 ? 'rgba(251,7,7,0.46)' :
            log.level === 'WARNING' ? 'rgba(255, 140, 0, 0.1)' :
            'transparent'
          }
            >
              {#each columns as column}
                <td
                        class="p-2 border-b"
                        style:min-width="{column.minWidth}px"
                >
                  {#if column.id === 'timestamp'}
                    {new Date(log[column.id]).toLocaleString()}
                  {:else if column.id === 'data'}
                    {JSON.stringify(log[column.id]?.fields || {})}
                  {:else}
                    {log[column.id]}
                  {/if}
                </td>
              {/each}
            </tr>
          {/each}
          </tbody>
        </table>


  </section>

</main>

<style lang="scss">

  .control {
    display: flex;
    justify-content: space-between;
    align-items: center;
    width: 100%;
    color: #ffffff;
    font-family: "Cascadia Mono", monospace;
    border-top: 1px solid #2b2b2b;
    border-bottom: 1px solid #2b2b2b;

    > div {
      display: flex;
      align-items: center;
      padding: 0.5rem 1rem;
    }

    .options {
      display: flex;
      gap: 1rem;
    }

    button {
      padding: 0.5rem 1rem;
      background: #3f8181;
      color: #ffffff;
      border: 1px solid #2b2b2b;
      border-radius: 0.25rem;
      cursor: pointer;
      font-family: "Cascadia Mono", monospace;
    }

    input {
      padding: 0.5rem;
      border: 1px solid #2b2b2b;
      border-radius: 0.25rem;
      font-family: "Cascadia Mono", monospace;
      background: #000000;
        color: #ffffff;
    }
  }

    main {
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        height: 100dvh;
        width: 100%;
      > section {
        display: flex;
        flex-direction: column;
        align-items: center;
        flex: 1;
        width: 100%;
        overflow: scroll;
        background: #000000;
      }


    }

    table {
      width: 100%;
      border-collapse: collapse;
      border-spacing: 0;
      border: 1px solid #2b2b2b;
      background: #000000;
      color: #ffffff;
      font-family: "Cascadia Mono", monospace;
    }

    th, td {
      padding: 0.2rem;
      border-bottom: 1px solid #2b2b2b;
      border-right: 1px solid #2b2b2b;
      text-align: left;
      font-size: 0.8rem;
    }

    th {
      background: #2b2b2b;
      padding: 0.5rem;
      color: #ffffff;
      text-align: left;
    }



</style>
