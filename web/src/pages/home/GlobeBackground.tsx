import { useEffect, useRef } from "react";
import { isLandAt } from "@/lib/globe/landMask";
import { subsolarPoint } from "@/lib/globe/solar";
import {
  getUtcOffset,
  knownZoneCoords,
  nearestTimeZone,
  timeZoneToLatLng,
  ZONE_ENTRIES,
} from "@/lib/globe/timezoneCoords";

/**
 * Decorative dot-matrix globe behind the home hero (D-10 / DESIGN.md
 * "ホーム・ヒーロー背景"). Renders the selected time zone as: the globe rotated so
 * that meridian faces front, a pulsing marker, a day/night terminator, and an
 * offset arc + label. Display-only — never used for correctness.
 *
 * Constraints honoured here: Three.js is lazy-imported (kept off other routes),
 * colours are read from CSS tokens (theme-aware, no dark-variant classes), motion stops for
 * `prefers-reduced-motion` and hidden tabs, and a missing WebGL context is a
 * silent no-op.
 */

const DEG = Math.PI / 180;

interface GlobeBackgroundProps {
  timeZone: string;
  /** Called when the user clicks a location on the globe (nearest zone). */
  onTimeZoneChange?: (timeZone: string) => void;
}

function hasWebGL(): boolean {
  try {
    const canvas = document.createElement("canvas");
    return (
      typeof WebGLRenderingContext !== "undefined" &&
      Boolean(canvas.getContext("webgl2") ?? canvas.getContext("webgl"))
    );
  } catch {
    return false;
  }
}

/** lat/lng (deg) → unit vector. lng 0 faces +Z (camera), east is +X. */
function toVec3(lat: number, lng: number): [number, number, number] {
  const phi = (90 - lat) * DEG;
  const theta = lng * DEG;
  const s = Math.sin(phi);
  return [s * Math.sin(theta), Math.cos(phi), s * Math.cos(theta)];
}

/**
 * True when this runtime's `Intl` resolves the zone id. Guards ids coming out
 * of the boundary dataset before they reach app state (picker / Temporal).
 */
function isUsableTimeZone(timeZone: string): boolean {
  try {
    new Intl.DateTimeFormat("en-US", { timeZone });
    return true;
  } catch {
    return false;
  }
}

export function GlobeBackground({ timeZone, onTimeZoneChange }: GlobeBackgroundProps) {
  const containerRef = useRef<HTMLDivElement>(null);
  const hitRef = useRef<HTMLDivElement>(null);
  const labelRef = useRef<HTMLDivElement>(null);
  const nameRef = useRef<HTMLSpanElement>(null);
  const offsetRef = useRef<HTMLSpanElement>(null);
  const hoverLabelRef = useRef<HTMLDivElement>(null);
  const hoverNameRef = useRef<HTMLSpanElement>(null);
  const hoverOffsetRef = useRef<HTMLSpanElement>(null);
  // Bridge live prop updates into the imperative Three.js scene without
  // rebuilding it on every time-zone change.
  const apiRef = useRef<((tz: string) => void) | null>(null);
  const tzRef = useRef(timeZone);
  const onSelectRef = useRef(onTimeZoneChange);

  useEffect(() => {
    tzRef.current = timeZone;
    onSelectRef.current = onTimeZoneChange;
    apiRef.current?.(timeZone);
  }, [timeZone, onTimeZoneChange]);

  // biome-ignore lint/correctness/useExhaustiveDependencies: scene is built once on mount; prop changes flow through apiRef.
  useEffect(() => {
    if (!hasWebGL()) return;

    let cancelled = false;
    let dispose = () => {};

    (async () => {
      // Both stay off other routes' bundles: three (renderer) and the
      // timezone-boundary-builder-derived lat/lng → IANA zone lookup that
      // drives the zone fill and click picking (display-only, D-10).
      const [THREE, { default: tzLookup }] = await Promise.all([
        import("three"),
        import("@photostructure/tz-lookup"),
      ]);
      const container = containerRef.current;
      if (cancelled || !container) return;

      const reducedMotion = window.matchMedia("(prefers-reduced-motion: reduce)");

      // --- colours from CSS tokens (rasterise the oklch value to read back RGB) ---
      const probe = document.createElement("canvas");
      probe.width = probe.height = 1;
      const probeCtx = probe.getContext("2d");
      const tokenColor = (name: string): InstanceType<typeof THREE.Color> => {
        const raw =
          getComputedStyle(document.documentElement).getPropertyValue(name).trim() || "#888";
        if (!probeCtx) return new THREE.Color(raw);
        probeCtx.fillStyle = "#000";
        probeCtx.fillStyle = raw;
        probeCtx.fillRect(0, 0, 1, 1);
        const [r, g, b] = probeCtx.getImageData(0, 0, 1, 1).data;
        return new THREE.Color(r / 255, g / 255, b / 255);
      };

      // --- renderer / scene / camera ---
      const renderer = new THREE.WebGLRenderer({ antialias: true, alpha: true });
      renderer.setPixelRatio(Math.min(2, window.devicePixelRatio));
      renderer.setClearColor(0x000000, 0);
      container.appendChild(renderer.domElement);

      const scene = new THREE.Scene();
      const camera = new THREE.PerspectiveCamera(32, 1, 0.1, 100);
      camera.position.set(0, 0, 4.6);

      const tilt = new THREE.Group();
      // vertical rotation follows the selected city's latitude (set per selection)
      scene.add(tilt);
      const spin = new THREE.Group();
      tilt.add(spin);

      // Curated zones' positions on the unit sphere — the always-visible
      // city markers with hover labels. The zone *fill* below no longer
      // derives from these points; it comes from real boundary data.
      const zoneVecs = ZONE_ENTRIES.map(([, c]) => toVec3(c.lat, c.lng));

      // Zone-fill index: every land dot is assigned its actual IANA zone via
      // tzLookup (timezone-boundary-builder data), so the highlight follows
      // real time-zone borders — Hokkaido/Okinawa land with Asia/Tokyo, all
      // of China with Asia/Shanghai, etc. Interned to a dense int per zone
      // name so the shader compares a single float (display-only, D-10).
      const zoneIndexByName = new Map<string, number>();
      const internZone = (zone: string): number => {
        let zi = zoneIndexByName.get(zone);
        if (zi === undefined) {
          zi = zoneIndexByName.size;
          zoneIndexByName.set(zone, zi);
        }
        return zi;
      };

      // --- dot cloud: Fibonacci sphere, land bright / ocean faint ---
      const N = 384000;
      const golden = Math.PI * (3 - Math.sqrt(5));
      const positions = new Float32Array(N * 3);
      const landAttr = new Float32Array(N);
      const zoneAttr = new Float32Array(N);
      for (let i = 0; i < N; i++) {
        const y = 1 - (i / (N - 1)) * 2;
        const r = Math.sqrt(Math.max(0, 1 - y * y));
        const th = i * golden;
        const x = Math.cos(th) * r;
        const z = Math.sin(th) * r;
        positions[i * 3] = x;
        positions[i * 3 + 1] = y;
        positions[i * 3 + 2] = z;
        const lat = Math.asin(y) / DEG;
        const lng = Math.atan2(x, z) / DEG;
        const land = isLandAt(lat, lng);
        landAttr[i] = land ? 1 : 0;
        if (land) zoneAttr[i] = internZone(tzLookup(lat, lng));
      }
      const geo = new THREE.BufferGeometry();
      geo.setAttribute("position", new THREE.BufferAttribute(positions, 3));
      geo.setAttribute("aLand", new THREE.BufferAttribute(landAttr, 1));
      geo.setAttribute("aZoneIndex", new THREE.BufferAttribute(zoneAttr, 1));

      const uniforms = {
        uSun: { value: new THREE.Vector3(...toVec3(0, 0)) },
        uColorLand: { value: tokenColor("--foreground") },
        uColorOcean: { value: tokenColor("--muted-foreground") },
        uAlphaLand: { value: 0.5 },
        uAlphaOcean: { value: 0.16 },
        uNightDim: { value: 0.4 },
        uSize: { value: 1.8 },
        uPixelRatio: { value: renderer.getPixelRatio() },
        // Highlight every land dot whose boundary-data zone matches the
        // selection — "this time zone's actual territory" (display-only, D-10).
        // -1 = nothing highlighted (zone with no land dots at this density).
        uColorZone: { value: tokenColor("--ixdtf-timezone") },
        uSelectedZoneIndex: { value: -1 },
      };
      const dotMaterial = new THREE.ShaderMaterial({
        uniforms,
        transparent: true,
        depthWrite: false,
        vertexShader: /* glsl */ `
          attribute float aLand;
          attribute float aZoneIndex;
          uniform vec3 uSun;
          uniform vec3 uColorLand;
          uniform vec3 uColorOcean;
          uniform float uAlphaLand;
          uniform float uAlphaOcean;
          uniform float uNightDim;
          uniform float uSize;
          uniform float uPixelRatio;
          uniform vec3 uColorZone;
          uniform float uSelectedZoneIndex;
          varying vec3 vColor;
          varying float vAlpha;
          void main() {
            vec4 mv = modelViewMatrix * vec4(position, 1.0);
            gl_PointSize = uSize * uPixelRatio * (8.0 / -mv.z);
            gl_Position = projectionMatrix * mv;
            float day = smoothstep(-0.15, 0.25, dot(normalize(position), uSun));
            float dim = mix(uNightDim, 1.0, day);
            vec3 base = mix(uColorOcean, uColorLand, aLand);
            float baseA = mix(uAlphaOcean, uAlphaLand, aLand) * dim;
            // land dots whose boundary-data zone matches the selection take the zone colour
            float inZone = (1.0 - step(0.5, abs(aZoneIndex - uSelectedZoneIndex))) * aLand;
            vColor = mix(base, uColorZone, inZone);
            vAlpha = mix(baseA, mix(0.55, 0.85, day), inZone);
          }
        `,
        fragmentShader: /* glsl */ `
          varying vec3 vColor;
          varying float vAlpha;
          void main() {
            vec2 c = gl_PointCoord - 0.5;
            float d = dot(c, c);
            if (d > 0.25) discard;
            float edge = smoothstep(0.25, 0.17, d);
            gl_FragColor = vec4(vColor, vAlpha * edge);
          }
        `,
      });
      const dots = new THREE.Points(geo, dotMaterial);
      spin.add(dots);

      // --- curated-city markers: always visible, hover shows the zone name ---
      const cityPositions = new Float32Array(zoneVecs.length * 3);
      zoneVecs.forEach(([vx, vy, vz], idx) => {
        cityPositions[idx * 3] = vx * 1.012;
        cityPositions[idx * 3 + 1] = vy * 1.012;
        cityPositions[idx * 3 + 2] = vz * 1.012;
      });
      const cityGeo = new THREE.BufferGeometry();
      cityGeo.setAttribute("position", new THREE.BufferAttribute(cityPositions, 3));
      const cityUniforms = {
        uColor: { value: tokenColor("--foreground") },
        uSize: { value: 4.5 },
        uPixelRatio: { value: renderer.getPixelRatio() },
      };
      const cityMaterial = new THREE.ShaderMaterial({
        uniforms: cityUniforms,
        transparent: true,
        depthWrite: false,
        vertexShader: /* glsl */ `
          uniform float uSize;
          uniform float uPixelRatio;
          void main() {
            vec4 mv = modelViewMatrix * vec4(position, 1.0);
            gl_PointSize = uSize * uPixelRatio * (8.0 / -mv.z);
            gl_Position = projectionMatrix * mv;
          }
        `,
        fragmentShader: /* glsl */ `
          uniform vec3 uColor;
          void main() {
            vec2 c = gl_PointCoord - 0.5;
            float d = dot(c, c);
            if (d > 0.25) discard;
            float edge = smoothstep(0.25, 0.17, d);
            gl_FragColor = vec4(uColor, 0.4 * edge);
          }
        `,
      });
      const cityPoints = new THREE.Points(cityGeo, cityMaterial);
      spin.add(cityPoints);

      // --- selection marker (ixdtf-timezone) + pulse ring ---
      const markerColor = tokenColor("--ixdtf-timezone");
      const markerDot = new THREE.Mesh(
        new THREE.SphereGeometry(0.022, 16, 16),
        new THREE.MeshBasicMaterial({ color: markerColor }),
      );
      spin.add(markerDot);
      const ringMaterial = new THREE.MeshBasicMaterial({
        color: markerColor,
        transparent: true,
        opacity: 0.7,
        side: THREE.DoubleSide,
        depthWrite: false,
      });
      const pulseRing = new THREE.Mesh(new THREE.RingGeometry(0.03, 0.045, 40), ringMaterial);
      spin.add(pulseRing);

      // --- offset meridian arc (ixdtf-offset) ---
      const arcMaterial = new THREE.LineBasicMaterial({
        color: tokenColor("--ixdtf-offset"),
        transparent: true,
        opacity: 0.5,
      });
      const arc = new THREE.Line(new THREE.BufferGeometry(), arcMaterial);
      spin.add(arc);

      // --- selection state, updated imperatively on time-zone change ---
      let targetRotY = 0;
      let targetTiltX = 0;
      let dirty = true;
      const updateSelection = (tz: string) => {
        const { lat, lng } = timeZoneToLatLng(tz);
        // nearest equivalent angle, so the ease never unwinds the full turns
        // accumulated by the idle spin (or by long drags)
        const baseRotY = -lng * DEG;
        targetRotY =
          baseRotY + Math.round((spin.rotation.y - baseRotY) / (2 * Math.PI)) * 2 * Math.PI;
        // rotate the city's latitude to the front (vertical) centre
        targetTiltX = Math.max(-1.25, Math.min(1.25, lat * DEG));
        const [mx, my, mz] = toVec3(lat, lng);
        markerDot.position.set(mx * 1.012, my * 1.012, mz * 1.012);
        pulseRing.position.copy(markerDot.position);
        pulseRing.quaternion.setFromUnitVectors(
          new THREE.Vector3(0, 0, 1),
          new THREE.Vector3(mx, my, mz).normalize(),
        );
        const pts: number[] = [];
        for (let a = -78; a <= 78; a += 3) {
          const [ax, ay, az] = toVec3(a, lng);
          pts.push(ax * 1.01, ay * 1.01, az * 1.01);
        }
        arc.geometry.dispose();
        arc.geometry = new THREE.BufferGeometry().setAttribute(
          "position",
          new THREE.Float32BufferAttribute(pts, 3),
        );
        // Resolve the selection to a fill index: exact zone-id match first;
        // for ids the boundary data spells differently (aliases), re-look-up
        // at the zone's real coordinates (curated city or tzdb zone.tab).
        // Zones with only a synthetic region/offset fallback get -1 (no
        // fill) rather than a guess from the wrong location.
        uniforms.uSelectedZoneIndex.value =
          zoneIndexByName.get(tz) ??
          (knownZoneCoords(tz) ? zoneIndexByName.get(tzLookup(lat, lng)) : undefined) ??
          -1;
        if (nameRef.current) nameRef.current.textContent = tz;
        if (offsetRef.current) offsetRef.current.textContent = getUtcOffset(tz).label;
        dirty = true;
      };
      updateSelection(tzRef.current);
      // On first mount, orient straight to the selection (no long spin-in).
      spin.rotation.y = targetRotY;
      tilt.rotation.x = targetTiltX;
      apiRef.current = updateSelection;

      // --- responsive sizing (offset the sphere off-centre on wide screens) ---
      let wide = false;
      const resize = () => {
        const w = container.clientWidth || window.innerWidth;
        const h = container.clientHeight || window.innerHeight;
        renderer.setSize(w, h, false);
        camera.aspect = w / h;
        wide = w >= 768;
        // Place the globe centre at a fixed viewport fraction so it stays framed
        // on the open right side at any width (content owns the narrow left).
        const span = 2 * camera.position.z * Math.tan((32 * DEG) / 2);
        const cx = wide ? (0.72 - 0.5) * (w / h) * span : 0;
        const cy = (0.5 - (wide ? 0.46 : 0.32)) * span;
        tilt.position.set(cx, cy, 0);
        camera.updateProjectionMatrix();
        dirty = true;
      };
      const ro = new ResizeObserver(resize);
      ro.observe(container);
      resize();

      // --- pointer interaction: drag to rotate, click to pick a zone ---
      const pickMesh = new THREE.Mesh(
        new THREE.SphereGeometry(1, 32, 24),
        new THREE.MeshBasicMaterial({ transparent: true, opacity: 0, depthWrite: false }),
      );
      spin.add(pickMesh);
      const raycaster = new THREE.Raycaster();
      raycaster.params.Points = { threshold: 0.025 };
      const ndc = new THREE.Vector2();
      const hit = hitRef.current;

      const pickPoint = (
        clientX: number,
        clientY: number,
        target: InstanceType<typeof THREE.Object3D> = pickMesh,
      ) => {
        const rect = renderer.domElement.getBoundingClientRect();
        ndc.set(
          ((clientX - rect.left) / rect.width) * 2 - 1,
          -((clientY - rect.top) / rect.height) * 2 + 1,
        );
        raycaster.setFromCamera(ndc, camera);
        return raycaster.intersectObject(target)[0];
      };

      let dragging = false;
      let moved = 0;
      let lastX = 0;
      let lastY = 0;

      const onHoverMove = (e: PointerEvent) => {
        if (dragging || !hit) return;
        const hoverLabel = hoverLabelRef.current;
        const cityHit = pickPoint(e.clientX, e.clientY, cityPoints);
        if (cityHit && typeof cityHit.index === "number") {
          if (hoverLabel) {
            const rect = renderer.domElement.getBoundingClientRect();
            const zone = ZONE_ENTRIES[cityHit.index][0];
            if (hoverNameRef.current) hoverNameRef.current.textContent = zone;
            if (hoverOffsetRef.current)
              hoverOffsetRef.current.textContent = getUtcOffset(zone).label;
            hoverLabel.style.transform = `translate3d(${e.clientX - rect.left}px, ${e.clientY - rect.top - 14}px, 0)`;
            hoverLabel.style.opacity = "1";
          }
          hit.style.cursor = "pointer";
          return;
        }
        if (hoverLabel) hoverLabel.style.opacity = "0";
        hit.style.cursor = wide && pickPoint(e.clientX, e.clientY) ? "grab" : "default";
      };
      const onPointerDown = (e: PointerEvent) => {
        if (!wide || !pickPoint(e.clientX, e.clientY)) return;
        dragging = true;
        moved = 0;
        lastX = e.clientX;
        lastY = e.clientY;
        if (hit) hit.style.cursor = "grabbing";
        hit?.setPointerCapture(e.pointerId);
      };
      const onPointerMove = (e: PointerEvent) => {
        if (!dragging) return;
        const dx = e.clientX - lastX;
        const dy = e.clientY - lastY;
        lastX = e.clientX;
        lastY = e.clientY;
        moved += Math.abs(dx) + Math.abs(dy);
        spin.rotation.y += dx * 0.005;
        targetRotY = spin.rotation.y; // keep the ease from fighting the drag
        tilt.rotation.x = Math.max(-1.25, Math.min(1.25, tilt.rotation.x + dy * 0.005));
        targetTiltX = tilt.rotation.x;
        dirty = true;
      };
      const onPointerUp = (e: PointerEvent) => {
        if (!dragging) return;
        dragging = false;
        hit?.releasePointerCapture(e.pointerId);
        if (hit) hit.style.cursor = "grab";
        if (moved >= 6) return; // a drag, not a click
        const p = pickPoint(e.clientX, e.clientY);
        if (!p) return;
        const local = spin.worldToLocal(p.point.clone()).normalize();
        const lat = Math.asin(Math.max(-1, Math.min(1, local.y))) / DEG;
        const lng = Math.atan2(local.x, local.z) / DEG;
        // Boundary-accurate pick; fall back to the curated nearest city only
        // when this runtime's Intl can't resolve the dataset's zone id.
        const zone = tzLookup(lat, lng);
        onSelectRef.current?.(isUsableTimeZone(zone) ? zone : nearestTimeZone(lat, lng));
      };
      hit?.addEventListener("pointerdown", onPointerDown);
      window.addEventListener("pointermove", onHoverMove);
      window.addEventListener("pointermove", onPointerMove);
      window.addEventListener("pointerup", onPointerUp);

      // --- theme / colour refresh ---
      const refreshColors = () => {
        uniforms.uColorLand.value = tokenColor("--foreground");
        uniforms.uColorOcean.value = tokenColor("--muted-foreground");
        uniforms.uColorZone.value = tokenColor("--ixdtf-timezone");
        cityUniforms.uColor.value = tokenColor("--foreground");
        (markerDot.material as InstanceType<typeof THREE.MeshBasicMaterial>).color =
          tokenColor("--ixdtf-timezone");
        ringMaterial.color = tokenColor("--ixdtf-timezone");
        arcMaterial.color = tokenColor("--ixdtf-offset");
        dirty = true;
      };
      const themeObserver = new MutationObserver(refreshColors);
      themeObserver.observe(document.documentElement, {
        attributes: true,
        attributeFilter: ["class"],
      });
      const schemeQuery = window.matchMedia("(prefers-color-scheme: dark)");
      schemeQuery.addEventListener("change", refreshColors);

      // --- render loop ---
      const project = new THREE.Vector3();
      let paused = document.hidden;
      let raf = 0;

      const render = () => {
        // day/night terminator from current time (local earth frame)
        const sun = subsolarPoint(new Date());
        uniforms.uSun.value.set(...toVec3(sun.lat, sun.lng));

        // offset label: project the marker to screen, hide when it faces away
        const label = labelRef.current;
        if (label) {
          markerDot.getWorldPosition(project);
          const facing = project.clone().normalize().z > 0.12;
          project.project(camera);
          if (facing && project.z < 1) {
            const w = container.clientWidth;
            const h = container.clientHeight;
            const x = (project.x * 0.5 + 0.5) * w;
            const y = (-project.y * 0.5 + 0.5) * h;
            label.style.transform = `translate3d(${x}px, ${y}px, 0)`;
            label.style.opacity = "1";
          } else {
            label.style.opacity = "0";
          }
        }
        renderer.render(scene, camera);
      };

      // idle spin: the same direction as Earth's real rotation, one
      // revolution ≈ 2 minutes — slow enough to stay "atmosphere" (D-10)
      const IDLE_SPIN = (2 * Math.PI) / 320;
      let lastFrame = performance.now();

      const loop = () => {
        raf = requestAnimationFrame(loop);
        if (paused) return;
        if (reducedMotion.matches) {
          if (!dirty) return;
          spin.rotation.y = targetRotY;
          tilt.rotation.x = targetTiltX;
          ringMaterial.opacity = 0.6;
          pulseRing.scale.setScalar(1);
          render();
          dirty = false;
          return;
        }
        const now = performance.now();
        // clamp so a backgrounded tab doesn't fast-forward the spin on return
        const dt = Math.min(0.1, (now - lastFrame) / 1000);
        lastFrame = now;
        // drift the target so the globe keeps turning; dragging pauses it
        if (!dragging) targetRotY += IDLE_SPIN * dt;
        // ease the selected meridian to the front (horizontal) and the selected
        // city's latitude to the vertical centre
        spin.rotation.y += (targetRotY - spin.rotation.y) * 0.06;
        tilt.rotation.x += (targetTiltX - tilt.rotation.x) * 0.06;
        // marker pulse
        const t = (now % 2600) / 2600;
        pulseRing.scale.setScalar(1 + t * 1.6);
        ringMaterial.opacity = 0.6 * (1 - t);
        render();
      };
      loop();

      const onVisibility = () => {
        paused = document.hidden;
        dirty = true;
      };
      document.addEventListener("visibilitychange", onVisibility);

      dispose = () => {
        cancelAnimationFrame(raf);
        document.removeEventListener("visibilitychange", onVisibility);
        hit?.removeEventListener("pointerdown", onPointerDown);
        window.removeEventListener("pointermove", onHoverMove);
        window.removeEventListener("pointermove", onPointerMove);
        window.removeEventListener("pointerup", onPointerUp);
        schemeQuery.removeEventListener("change", refreshColors);
        themeObserver.disconnect();
        ro.disconnect();
        pickMesh.geometry.dispose();
        (pickMesh.material as InstanceType<typeof THREE.Material>).dispose();
        geo.dispose();
        dotMaterial.dispose();
        cityGeo.dispose();
        cityMaterial.dispose();
        markerDot.geometry.dispose();
        (markerDot.material as InstanceType<typeof THREE.Material>).dispose();
        pulseRing.geometry.dispose();
        ringMaterial.dispose();
        arc.geometry.dispose();
        arcMaterial.dispose();
        renderer.dispose();
        renderer.domElement.remove();
      };
    })();

    return () => {
      cancelled = true;
      apiRef.current = null;
      dispose();
    };
  }, []);

  return (
    <>
      {/* md 未満はコンテンツが全幅で地球儀と重なるため、全体を減光して
          文字の可読性を守る (DESIGN.md「文字が主役、地球儀は大気」)。 */}
      <div
        ref={containerRef}
        aria-hidden="true"
        className="pointer-events-none fixed inset-0 z-0 overflow-hidden opacity-40 md:opacity-100 [&>canvas]:block"
      >
        <div
          ref={labelRef}
          className="absolute top-0 left-0 flex flex-col whitespace-nowrap rounded-sm bg-background/80 px-1.5 py-0.5 font-medium font-mono text-xs leading-tight tabular-nums opacity-0 transition-opacity duration-300"
          style={{ willChange: "transform" }}
        >
          <span ref={nameRef} className="text-ixdtf-timezone" />
          <span ref={offsetRef} className="text-ixdtf-offset" />
        </div>
        {/* Hover label: identical styling to the selection label above. */}
        <div
          ref={hoverLabelRef}
          className="absolute top-0 left-0 flex flex-col whitespace-nowrap rounded-sm bg-background/80 px-1.5 py-0.5 font-medium font-mono text-xs leading-tight tabular-nums opacity-0 transition-opacity duration-300"
          style={{ willChange: "transform" }}
        >
          <span ref={hoverNameRef} className="text-ixdtf-timezone" />
          <span ref={hoverOffsetRef} className="text-ixdtf-offset" />
        </div>
      </div>
      {/* Interaction layer: sits below the content column (z-10) so cards always
          win, above the globe canvas so drags/clicks land on the sphere (D-10). */}
      <div ref={hitRef} aria-hidden="true" className="fixed inset-0 z-[5]" />
    </>
  );
}
